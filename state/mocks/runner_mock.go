// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/ij/command -i Runner -d mocks -f

package mocks

import (
	"context"
	logging "github.com/efritz/ij/logging"
	"sync"
)

type MockRunner struct {
	RunFunc          func(context.Context, []string, *logging.Prefix) error
	histRun          []RunnerRunParamSet
	RunForOutputFunc func(context.Context, []string) (string, string, error)
	histRunForOutput []RunnerRunForOutputParamSet
	mutex            sync.RWMutex
}
type RunnerRunParamSet struct {
	Arg0 context.Context
	Arg1 []string
	Arg2 *logging.Prefix
}
type RunnerRunForOutputParamSet struct {
	Arg0 context.Context
	Arg1 []string
}

func NewMockRunner() *MockRunner {
	m := &MockRunner{}
	m.RunFunc = m.defaultRunFunc
	m.RunForOutputFunc = m.defaultRunForOutputFunc
	return m
}
func (m *MockRunner) Run(v0 context.Context, v1 []string, v2 *logging.Prefix) error {
	m.mutex.Lock()
	m.histRun = append(m.histRun, RunnerRunParamSet{v0, v1, v2})
	m.mutex.Unlock()
	return m.RunFunc(v0, v1, v2)
}
func (m *MockRunner) RunFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histRun)
}
func (m *MockRunner) RunFuncCallParams() []RunnerRunParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histRun
}

func (m *MockRunner) RunForOutput(v0 context.Context, v1 []string) (string, string, error) {
	m.mutex.Lock()
	m.histRunForOutput = append(m.histRunForOutput, RunnerRunForOutputParamSet{v0, v1})
	m.mutex.Unlock()
	return m.RunForOutputFunc(v0, v1)
}
func (m *MockRunner) RunForOutputFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histRunForOutput)
}
func (m *MockRunner) RunForOutputFuncCallParams() []RunnerRunForOutputParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histRunForOutput
}

func (m *MockRunner) defaultRunFunc(v0 context.Context, v1 []string, v2 *logging.Prefix) error {
	return nil
}
func (m *MockRunner) defaultRunForOutputFunc(v0 context.Context, v1 []string) (string, string, error) {
	return "", "", nil
}
