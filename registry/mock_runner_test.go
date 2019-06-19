// Code generated by github.com/efritz/go-mockgen 0.1.0; DO NOT EDIT.
// This file was generated by robots at
// 2019-06-19T11:52:20-05:00
// using the command
// $ go-mockgen -f github.com/ij-build/ij/command -i Runner -o mock_runner_test.go

package registry

import (
	"context"
	command "github.com/ij-build/ij/command"
	logging "github.com/ij-build/ij/logging"
	"io"
	"sync"
)

// MockRunner is a mock implementation of the Runner interface (from the
// package github.com/ij-build/ij/command) used for unit testing.
type MockRunner struct {
	// RunFunc is an instance of a mock function object controlling the
	// behavior of the method Run.
	RunFunc *RunnerRunFunc
	// RunForOutputFunc is an instance of a mock function object controlling
	// the behavior of the method RunForOutput.
	RunForOutputFunc *RunnerRunForOutputFunc
}

// NewMockRunner creates a new mock of the Runner interface. All methods
// return zero values for all results, unless overwritten.
func NewMockRunner() *MockRunner {
	return &MockRunner{
		RunFunc: &RunnerRunFunc{
			defaultHook: func(context.Context, []string, io.ReadCloser, *logging.Prefix) error {
				return nil
			},
		},
		RunForOutputFunc: &RunnerRunForOutputFunc{
			defaultHook: func(context.Context, []string, io.ReadCloser) (string, string, error) {
				return "", "", nil
			},
		},
	}
}

// NewMockRunnerFrom creates a new mock of the MockRunner interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockRunnerFrom(i command.Runner) *MockRunner {
	return &MockRunner{
		RunFunc: &RunnerRunFunc{
			defaultHook: i.Run,
		},
		RunForOutputFunc: &RunnerRunForOutputFunc{
			defaultHook: i.RunForOutput,
		},
	}
}

// RunnerRunFunc describes the behavior when the Run method of the parent
// MockRunner instance is invoked.
type RunnerRunFunc struct {
	defaultHook func(context.Context, []string, io.ReadCloser, *logging.Prefix) error
	hooks       []func(context.Context, []string, io.ReadCloser, *logging.Prefix) error
	history     []RunnerRunFuncCall
	mutex       sync.Mutex
}

// Run delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockRunner) Run(v0 context.Context, v1 []string, v2 io.ReadCloser, v3 *logging.Prefix) error {
	r0 := m.RunFunc.nextHook()(v0, v1, v2, v3)
	m.RunFunc.appendCall(RunnerRunFuncCall{v0, v1, v2, v3, r0})
	return r0
}

// SetDefaultHook sets function that is called when the Run method of the
// parent MockRunner instance is invoked and the hook queue is empty.
func (f *RunnerRunFunc) SetDefaultHook(hook func(context.Context, []string, io.ReadCloser, *logging.Prefix) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Run method of the parent MockRunner instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *RunnerRunFunc) PushHook(hook func(context.Context, []string, io.ReadCloser, *logging.Prefix) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *RunnerRunFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, []string, io.ReadCloser, *logging.Prefix) error {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *RunnerRunFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context, []string, io.ReadCloser, *logging.Prefix) error {
		return r0
	})
}

func (f *RunnerRunFunc) nextHook() func(context.Context, []string, io.ReadCloser, *logging.Prefix) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *RunnerRunFunc) appendCall(r0 RunnerRunFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of RunnerRunFuncCall objects describing the
// invocations of this function.
func (f *RunnerRunFunc) History() []RunnerRunFuncCall {
	f.mutex.Lock()
	history := make([]RunnerRunFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// RunnerRunFuncCall is an object that describes an invocation of method Run
// on an instance of MockRunner.
type RunnerRunFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 []string
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 io.ReadCloser
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 *logging.Prefix
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c RunnerRunFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c RunnerRunFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// RunnerRunForOutputFunc describes the behavior when the RunForOutput
// method of the parent MockRunner instance is invoked.
type RunnerRunForOutputFunc struct {
	defaultHook func(context.Context, []string, io.ReadCloser) (string, string, error)
	hooks       []func(context.Context, []string, io.ReadCloser) (string, string, error)
	history     []RunnerRunForOutputFuncCall
	mutex       sync.Mutex
}

// RunForOutput delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockRunner) RunForOutput(v0 context.Context, v1 []string, v2 io.ReadCloser) (string, string, error) {
	r0, r1, r2 := m.RunForOutputFunc.nextHook()(v0, v1, v2)
	m.RunForOutputFunc.appendCall(RunnerRunForOutputFuncCall{v0, v1, v2, r0, r1, r2})
	return r0, r1, r2
}

// SetDefaultHook sets function that is called when the RunForOutput method
// of the parent MockRunner instance is invoked and the hook queue is empty.
func (f *RunnerRunForOutputFunc) SetDefaultHook(hook func(context.Context, []string, io.ReadCloser) (string, string, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// RunForOutput method of the parent MockRunner instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *RunnerRunForOutputFunc) PushHook(hook func(context.Context, []string, io.ReadCloser) (string, string, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *RunnerRunForOutputFunc) SetDefaultReturn(r0 string, r1 string, r2 error) {
	f.SetDefaultHook(func(context.Context, []string, io.ReadCloser) (string, string, error) {
		return r0, r1, r2
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *RunnerRunForOutputFunc) PushReturn(r0 string, r1 string, r2 error) {
	f.PushHook(func(context.Context, []string, io.ReadCloser) (string, string, error) {
		return r0, r1, r2
	})
}

func (f *RunnerRunForOutputFunc) nextHook() func(context.Context, []string, io.ReadCloser) (string, string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *RunnerRunForOutputFunc) appendCall(r0 RunnerRunForOutputFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of RunnerRunForOutputFuncCall objects
// describing the invocations of this function.
func (f *RunnerRunForOutputFunc) History() []RunnerRunForOutputFuncCall {
	f.mutex.Lock()
	history := make([]RunnerRunForOutputFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// RunnerRunForOutputFuncCall is an object that describes an invocation of
// method RunForOutput on an instance of MockRunner.
type RunnerRunForOutputFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 []string
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 io.ReadCloser
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 string
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c RunnerRunForOutputFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c RunnerRunForOutputFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2}
}
