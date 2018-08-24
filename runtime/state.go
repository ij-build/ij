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
	runID               string
	exportedEnv         []string
	envMutex            sync.RWMutex
	cpuShares           string
	ctx                 context.Context
	enableSSHAgent      bool
	env                 []string
	forceSequential     bool
	healthcheckInterval time.Duration
	memory              string
	cancel              func()
	once                sync.Once
	cleanup             *Cleanup
	containerStopper    *ContainerList
	logger              logging.Logger
	logProcessor        logging.Processor
	network             *Network
	networkDisconnector *ContainerList
	scratch             *ScratchSpace
}

func NewState(
	config *config.Config,
	plans []string,
	colorize bool,
	cpuShares string,
	enableSSHAgent bool,
	env []string,
	forceSequential bool,
	healthcheckInterval time.Duration,
	keepWorkspace bool,
	memory string,
	planTimeout time.Duration,
	verbose bool,
) (s *State, err error) {
	ctx, cancel := makeContext(planTimeout)

	s = &State{
		config:              config,
		plans:               plans,
		env:                 env,
		cpuShares:           cpuShares,
		enableSSHAgent:      enableSSHAgent,
		forceSequential:     forceSequential,
		healthcheckInterval: healthcheckInterval,
		memory:              memory,
		ctx:                 ctx,
		cancel:              cancel,
		cleanup:             NewCleanup(),
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
	// Generate a scratch directory

	s.scratch = NewScratchSpace(s.runID, keepWorkspace)

	if err = s.scratch.Setup(); err != nil {
		logging.EmergencyLog(
			"error: failed to create scratch directory: %s",
			err.Error(),
		)

		return
	}

	s.cleanup.Register(func() {
		if err := s.scratch.Prune(); err != nil {
			logging.EmergencyLog(
				"error: failed to clean up scratch directory: %s",
				err.Error(),
			)
		}
	})

	// If any of the remaining initialization fails, we
	// don't want to keep a scratch directory around so
	// we destory it at the end of the function on a
	// non-nil error return.

	defer func() {
		if err == nil {
			return
		}

		if err := s.scratch.Teardown(); err != nil {
			logging.EmergencyLog(
				"error: failed to teardown scratch directory: %s",
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

func (s *State) ExportEnv(line string) {
	s.envMutex.Lock()
	s.exportedEnv = append(s.exportedEnv, line)
	s.envMutex.Unlock()
}

func (s *State) GetExportedEnv() []string {
	s.envMutex.RLock()
	defer s.envMutex.RUnlock()

	env := []string{}
	for _, line := range s.exportedEnv {
		env = append(env, line)
	}

	return s.exportedEnv
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
