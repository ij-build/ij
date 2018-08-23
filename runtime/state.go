package runtime

import (
	"context"
	"sync"
	"time"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/util"
)

type State struct {
	config              *config.Config
	plans               []string
	env                 []string
	forceSequential     bool
	enableSSHAgent      bool
	healthcheckInterval time.Duration
	cpuShares           string
	memory              string
	cleanup             *Cleanup
	ctx                 context.Context
	cancel              func()
	once                sync.Once
	runID               string
	scratch             *ScratchSpace
	logProcessor        logging.Processor
	logger              logging.Logger
	containerStopper    *ContainerList
	networkDisconnector *ContainerList
	network             *Network
}

func NewState(
	config *config.Config,
	plans []string,
	env []string,
	verbose bool,
	colorize bool,
	forceSequential bool,
	enableSSHAgent bool,
	healthcheckInterval time.Duration,
	cpuShares string,
	memory string,
	planTimeout time.Duration,
) (s *State, err error) {
	ctx, cancel := makeContext(planTimeout)

	s = &State{
		config:              config,
		plans:               plans,
		env:                 env,
		forceSequential:     forceSequential,
		enableSSHAgent:      enableSSHAgent,
		healthcheckInterval: healthcheckInterval,
		cpuShares:           cpuShares,
		memory:              memory,
		cleanup:             NewCleanup(),
		ctx:                 ctx,
		cancel:              cancel,
	}

	//
	// Generate a unique Run ID

	if s.runID, err = util.MakeID(); err != nil {
		logging.EmergencyLog(
			"error: failed to generate run id: %s",
			err.Error(),
		)

		return
	}

	//
	// Generate a build directory

	s.scratch = NewScratchSpace(s.runID)

	if err = s.scratch.Setup(); err != nil {
		logging.EmergencyLog(
			"error: failed to create build directory: %s",
			err.Error(),
		)

		return
	}

	// If any of the remaining initialization fails, we
	// don't want to keep a build directory around so
	// we destory it at the end of the function on a
	// non-nil error return.

	defer func() {
		if err == nil {
			return
		}

		if err := s.scratch.Teardown(); err != nil {
			logging.EmergencyLog(
				"error: failed to teardown build directory: %s",
				err.Error(),
			)
		}
	}()

	//
	// Setup Logging and start log processor

	s.logProcessor = logging.NewProcessor(verbose, colorize)
	s.logProcessor.Start()
	s.cleanup.Register(s.logProcessor.Shutdown)

	//
	// Create Base Logger

	outfile, errfile, err := s.scratch.MakeLogFiles("ij")
	if err != nil {
		logging.EmergencyLog(
			"error: failed to create log files: %s",
			err.Error(),
		)

		return nil, err
	}

	s.logger = s.logProcessor.Logger(
		outfile,
		errfile,
		true,
	)

	//
	// Create Container Lists

	s.containerStopper = NewContainerStopper(
		s.logger,
	)

	s.networkDisconnector = NewNetworkDisconnector(
		s.runID,
		s.logger,
	)

	s.cleanup.Register(s.containerStopper.Execute)
	s.cleanup.Register(s.networkDisconnector.Execute)

	//
	// Create Network

	if s.network, err = NewNetwork(s.ctx, s.runID, s.logger); err != nil {
		s.ReportError(
			nil,
			"error: failed to create network: %s",
			err.Error(),
		)

		return
	}

	s.cleanup.Register(s.network.Teardown)
	return
}

func (s *State) ReportError(
	prefix *logging.Prefix,
	format string,
	args ...interface{},
) {
	select {
	case <-s.ctx.Done():
		return
	default:
	}

	s.logger.Error(prefix, format, args...)
}

func makeContext(timeout time.Duration) (context.Context, func()) {
	if timeout == 0 {
		return context.WithCancel(context.Background())
	}

	return context.WithTimeout(context.Background(), timeout)
}
