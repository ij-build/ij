package runtime

import (
	"context"
	"sync"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/util"
)

type State struct {
	config              *config.Config
	plans               []string
	env                 []string
	forceSequential     bool
	cleanup             *Cleanup
	ctx                 context.Context
	cancel              func()
	once                sync.Once
	runID               string
	buildDir            *BuildDir
	logProcessor        logging.Processor
	logger              logging.Logger
	containerStopper    *ContainerList
	networkDisconnector *ContainerList
	workspace           *Workspace
	network             *Network
}

func NewState(
	config *config.Config,
	plans []string,
	env []string,
	verbose bool,
	colorize bool,
	forceSequential bool,
) (s *State, err error) {
	ctx, cancel := context.WithCancel(context.Background())

	s = &State{
		config:          config,
		plans:           plans,
		env:             env,
		forceSequential: forceSequential,
		cleanup:         NewCleanup(),
		ctx:             ctx,
		cancel:          cancel,
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

	s.buildDir = NewBuildDir(s.runID)

	if err = s.buildDir.Setup(); err != nil {
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

		if err := s.buildDir.Teardown(); err != nil {
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

	outfile, errfile, err := s.buildDir.MakeLogFiles("ij")
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
	// Create Workspace

	if s.workspace, err = NewWorkspace(s.ctx, s.runID, s.logger); err != nil {
		s.ReportError(
			nil,
			"error: failed to create workspace volume: %s",
			err.Error(),
		)

		return
	}

	s.cleanup.Register(s.workspace.Teardown)

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
