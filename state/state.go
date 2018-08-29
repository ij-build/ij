package state

import (
	"context"
	"sync"
	"time"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/util"
)

type State struct {
	Config              *config.Config
	Plans               []string
	RunID               string
	exportedEnv         []string
	envMutex            sync.RWMutex
	CPUShares           string
	Context             context.Context
	EnableSSHAgent      bool
	Env                 []string
	ForceSequential     bool
	HealthcheckInterval time.Duration
	Memory              string
	Cancel              func()
	Once                sync.Once
	Cleanup             *Cleanup
	ContainerStopper    *ContainerList
	Logger              logging.Logger
	LogProcessor        logging.Processor
	NetworkDisconnector *ContainerList
	RegistryList        *RegistryList
	Scratch             *ScratchSpace
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
	login bool,
	memory string,
	planTimeout time.Duration,
	verbose bool,
) (s *State, err error) {
	ctx, cancel := makeContext(planTimeout)

	s = &State{
		Config:              config,
		Plans:               plans,
		Env:                 env,
		CPUShares:           cpuShares,
		EnableSSHAgent:      enableSSHAgent,
		ForceSequential:     forceSequential,
		HealthcheckInterval: healthcheckInterval,
		Memory:              memory,
		Context:             ctx,
		Cancel:              cancel,
		Cleanup:             NewCleanup(),
	}

	//
	// Generate a unique Run ID

	if s.RunID, err = util.MakeID(); err != nil {
		logging.EmergencyLog(
			"error: failed to generate run id: %s",
			err.Error(),
		)

		return
	}

	//
	// Generate a scratch directory

	s.Scratch = NewScratchSpace(s.RunID, keepWorkspace)

	if err = s.Scratch.Setup(); err != nil {
		logging.EmergencyLog(
			"error: failed to create scratch directory: %s",
			err.Error(),
		)

		return
	}

	s.Cleanup.Register(func() {
		if err := s.Scratch.Prune(); err != nil {
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

		if err := s.Scratch.Teardown(); err != nil {
			logging.EmergencyLog(
				"error: failed to teardown scratch directory: %s",
				err.Error(),
			)
		}
	}()

	//
	// Setup Logging and start log processor

	s.LogProcessor = logging.NewProcessor(verbose, colorize)
	s.LogProcessor.Start()
	s.Cleanup.Register(s.LogProcessor.Shutdown)

	//
	// Create Base Logger

	outfile, errfile, err := s.Scratch.MakeLogFiles("ij")
	if err != nil {
		logging.EmergencyLog(
			"error: failed to create log files: %s",
			err.Error(),
		)

		return nil, err
	}

	s.Logger = s.LogProcessor.Logger(
		outfile,
		errfile,
		true,
	)

	//
	// Create Container Lists

	s.ContainerStopper = NewContainerStopper(
		s.Logger,
	)

	s.NetworkDisconnector = NewNetworkDisconnector(
		s.RunID,
		s.Logger,
	)

	s.Cleanup.Register(s.ContainerStopper.Execute)
	s.Cleanup.Register(s.NetworkDisconnector.Execute)

	//
	// Login to Registries

	if login {
		registryEnv := environment.Merge(
			environment.New(s.Config.Environment),
			environment.New(s.Env),
		)

		registryList, registryErr := NewRegistryList(
			s.Context,
			s.Logger,
			registryEnv,
			s.Config.Registries,
		)

		if registryErr != nil {
			s.ReportError(
				nil,
				"error: failed to log into registries: %s",
				err.Error(),
			)

			err = registryErr
			return
		}

		s.Cleanup.Register(registryList.Teardown)
	}

	//
	// Create Network

	network, networkErr := NewNetwork(s.Context, s.RunID, s.Logger)
	if networkErr != nil {
		s.ReportError(
			nil,
			"error: failed to create network: %s",
			err.Error(),
		)

		err = networkErr
		return
	}

	s.Cleanup.Register(network.Teardown)
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
	case <-s.Context.Done():
		return
	default:
	}

	s.Logger.Error(prefix, format, args...)
}

func makeContext(timeout time.Duration) (context.Context, func()) {
	if timeout == 0 {
		return context.WithCancel(context.Background())
	}

	return context.WithTimeout(context.Background(), timeout)
}
