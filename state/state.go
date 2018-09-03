package state

import (
	"context"
	"sync"
	"time"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/registry"
	"github.com/efritz/ij/util"
)

type State struct {
	Config              *config.Config
	RunID               string
	exportedEnv         []string
	envMutex            sync.RWMutex
	CPUShares           string
	Context             context.Context
	EnableSSHAgent      bool
	Env                 []string
	Memory              string
	Cancel              func()
	Once                sync.Once
	Cleanup             *Cleanup
	ContainerStopper    *ContainerList
	Logger              logging.Logger
	LogProcessor        logging.Processor
	NetworkDisconnector *ContainerList
	RegistrySet         *registry.RegistrySet
	Scratch             *ScratchSpace
}

func NewState(
	config *config.Config,
	colorize bool,
	cpuShares string,
	enableSSHAgent bool,
	env []string,
	keepWorkspace bool,
	login bool,
	memory string,
	planTimeout time.Duration,
	verbose bool,
) (s *State, err error) {
	ctx, cancel := makeContext(planTimeout)

	s = &State{
		Config:         config,
		Env:            env,
		CPUShares:      cpuShares,
		EnableSSHAgent: enableSSHAgent,
		Memory:         memory,
		Context:        ctx,
		Cancel:         cancel,
		Cleanup:        NewCleanup(),
	}

	if err = s.setupRunID(); err != nil {
		return
	}

	if err = s.setupScratch(keepWorkspace); err != nil {
		return
	}

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

	if err = s.setupLogger(verbose, colorize); err != nil {
		return
	}

	if err = s.setupContainerLists(); err != nil {
		return
	}

	if err = s.setupRegistries(login); err != nil {
		return
	}

	if err = s.setupNetwork(); err != nil {
		return
	}

	return
}

func (s *State) setupRunID() error {
	id, err := util.MakeID()
	if err != nil {
		logging.EmergencyLog(
			"error: failed to generate run id: %s",
			err.Error(),
		)

		return err
	}

	s.RunID = id
	return nil
}

func (s *State) setupScratch(keepWorkspace bool) error {
	s.Scratch = NewScratchSpace(s.RunID, keepWorkspace)

	if err := s.Scratch.Setup(); err != nil {
		logging.EmergencyLog(
			"error: failed to create scratch directory: %s",
			err.Error(),
		)

		return err
	}

	s.Cleanup.Register(func() {
		if err := s.Scratch.Prune(); err != nil {
			logging.EmergencyLog(
				"error: failed to clean up scratch directory: %s",
				err.Error(),
			)
		}
	})

	return nil
}

func (s *State) setupLogger(verbose, colorize bool) error {
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

		return err
	}

	s.Logger = s.LogProcessor.Logger(
		outfile,
		errfile,
		true,
	)

	return nil
}

func (s *State) setupContainerLists() error {
	s.ContainerStopper = NewContainerStopper(
		s.Logger,
	)

	s.NetworkDisconnector = NewNetworkDisconnector(
		s.RunID,
		s.Logger,
	)

	s.Cleanup.Register(s.ContainerStopper.Execute)
	s.Cleanup.Register(s.NetworkDisconnector.Execute)
	return nil
}

func (s *State) setupRegistries(login bool) error {
	if !login {
		return nil
	}

	registryEnv := environment.Merge(
		environment.New(s.Config.Environment),
		environment.New(s.Env),
	)

	registrySet, err := registry.NewRegistrySet(
		s.Context,
		s.Logger,
		registryEnv,
		s.Config.Registries,
	)

	if err != nil {
		s.ReportError(
			nil,
			"error: failed to resolve registries: %s",
			err.Error(),
		)

		return err
	}

	if err = registrySet.Login(); err != nil {
		s.ReportError(
			nil,
			"error: failed to log into registries: %s",
			err.Error(),
		)

		return err
	}

	s.Cleanup.Register(registrySet.Logout)
	return nil
}

func (s *State) setupNetwork() error {
	network, err := NewNetwork(s.Context, s.RunID, s.Logger)
	if err != nil {
		s.ReportError(
			nil,
			"error: failed to create network: %s",
			err.Error(),
		)

		return err
	}

	s.Cleanup.Register(network.Teardown)
	return nil
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

func (s *State) BuildEnv(envs ...environment.Environment) environment.Environment {
	return environment.Merge(
		environment.New(s.Config.Environment),
		environment.Merge(envs...),
		environment.New(s.GetExportedEnv()),
		environment.New(s.Env),
	)
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
