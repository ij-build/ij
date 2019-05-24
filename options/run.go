package options

import (
	"context"
	"time"
)

type RunOptions struct {
	Plans                   []string
	CPUShares               string
	ForceSequential         bool
	HealthcheckInterval     time.Duration
	KeepWorkspace           bool
	Login                   bool
	Memory                  string
	PlanTimeout             time.Duration
	SSHIdentities           []string
	EnableContainerSSHAgent bool
	Context                 context.Context
}
