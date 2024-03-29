// Code generated by github.com/efritz/go-mockgen 0.1.0; DO NOT EDIT.
// This file was generated by robots at
// 2019-06-19T11:52:17-05:00
// using the command
// $ go-mockgen -f github.com/ij-build/ij/logging -i Logger -o mock_logger_test.go

package command

import (
	logging "github.com/ij-build/ij/logging"
	"sync"
)

// MockLogger is a mock implementation of the Logger interface (from the
// package github.com/ij-build/ij/logging) used for unit testing.
type MockLogger struct {
	// DebugFunc is an instance of a mock function object controlling the
	// behavior of the method Debug.
	DebugFunc *LoggerDebugFunc
	// ErrorFunc is an instance of a mock function object controlling the
	// behavior of the method Error.
	ErrorFunc *LoggerErrorFunc
	// InfoFunc is an instance of a mock function object controlling the
	// behavior of the method Info.
	InfoFunc *LoggerInfoFunc
	// WarnFunc is an instance of a mock function object controlling the
	// behavior of the method Warn.
	WarnFunc *LoggerWarnFunc
}

// NewMockLogger creates a new mock of the Logger interface. All methods
// return zero values for all results, unless overwritten.
func NewMockLogger() *MockLogger {
	return &MockLogger{
		DebugFunc: &LoggerDebugFunc{
			defaultHook: func(*logging.Prefix, string, ...interface{}) {
				return
			},
		},
		ErrorFunc: &LoggerErrorFunc{
			defaultHook: func(*logging.Prefix, string, ...interface{}) {
				return
			},
		},
		InfoFunc: &LoggerInfoFunc{
			defaultHook: func(*logging.Prefix, string, ...interface{}) {
				return
			},
		},
		WarnFunc: &LoggerWarnFunc{
			defaultHook: func(*logging.Prefix, string, ...interface{}) {
				return
			},
		},
	}
}

// NewMockLoggerFrom creates a new mock of the MockLogger interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockLoggerFrom(i logging.Logger) *MockLogger {
	return &MockLogger{
		DebugFunc: &LoggerDebugFunc{
			defaultHook: i.Debug,
		},
		ErrorFunc: &LoggerErrorFunc{
			defaultHook: i.Error,
		},
		InfoFunc: &LoggerInfoFunc{
			defaultHook: i.Info,
		},
		WarnFunc: &LoggerWarnFunc{
			defaultHook: i.Warn,
		},
	}
}

// LoggerDebugFunc describes the behavior when the Debug method of the
// parent MockLogger instance is invoked.
type LoggerDebugFunc struct {
	defaultHook func(*logging.Prefix, string, ...interface{})
	hooks       []func(*logging.Prefix, string, ...interface{})
	history     []LoggerDebugFuncCall
	mutex       sync.Mutex
}

// Debug delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockLogger) Debug(v0 *logging.Prefix, v1 string, v2 ...interface{}) {
	m.DebugFunc.nextHook()(v0, v1, v2...)
	m.DebugFunc.appendCall(LoggerDebugFuncCall{v0, v1, v2})
	return
}

// SetDefaultHook sets function that is called when the Debug method of the
// parent MockLogger instance is invoked and the hook queue is empty.
func (f *LoggerDebugFunc) SetDefaultHook(hook func(*logging.Prefix, string, ...interface{})) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Debug method of the parent MockLogger instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *LoggerDebugFunc) PushHook(hook func(*logging.Prefix, string, ...interface{})) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *LoggerDebugFunc) SetDefaultReturn() {
	f.SetDefaultHook(func(*logging.Prefix, string, ...interface{}) {
		return
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *LoggerDebugFunc) PushReturn() {
	f.PushHook(func(*logging.Prefix, string, ...interface{}) {
		return
	})
}

func (f *LoggerDebugFunc) nextHook() func(*logging.Prefix, string, ...interface{}) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *LoggerDebugFunc) appendCall(r0 LoggerDebugFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of LoggerDebugFuncCall objects describing the
// invocations of this function.
func (f *LoggerDebugFunc) History() []LoggerDebugFuncCall {
	f.mutex.Lock()
	history := make([]LoggerDebugFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// LoggerDebugFuncCall is an object that describes an invocation of method
// Debug on an instance of MockLogger.
type LoggerDebugFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 *logging.Prefix
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 string
	// Arg2 is a slice containing the values of the variadic arguments
	// passed to this method invocation.
	Arg2 []interface{}
}

// Args returns an interface slice containing the arguments of this
// invocation. The variadic slice argument is flattened in this array such
// that one positional argument and three variadic arguments would result in
// a slice of four, not two.
func (c LoggerDebugFuncCall) Args() []interface{} {
	trailing := []interface{}{}
	for _, val := range c.Arg2 {
		trailing = append(trailing, val)
	}

	return append([]interface{}{c.Arg0, c.Arg1}, trailing...)
}

// Results returns an interface slice containing the results of this
// invocation.
func (c LoggerDebugFuncCall) Results() []interface{} {
	return []interface{}{}
}

// LoggerErrorFunc describes the behavior when the Error method of the
// parent MockLogger instance is invoked.
type LoggerErrorFunc struct {
	defaultHook func(*logging.Prefix, string, ...interface{})
	hooks       []func(*logging.Prefix, string, ...interface{})
	history     []LoggerErrorFuncCall
	mutex       sync.Mutex
}

// Error delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockLogger) Error(v0 *logging.Prefix, v1 string, v2 ...interface{}) {
	m.ErrorFunc.nextHook()(v0, v1, v2...)
	m.ErrorFunc.appendCall(LoggerErrorFuncCall{v0, v1, v2})
	return
}

// SetDefaultHook sets function that is called when the Error method of the
// parent MockLogger instance is invoked and the hook queue is empty.
func (f *LoggerErrorFunc) SetDefaultHook(hook func(*logging.Prefix, string, ...interface{})) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Error method of the parent MockLogger instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *LoggerErrorFunc) PushHook(hook func(*logging.Prefix, string, ...interface{})) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *LoggerErrorFunc) SetDefaultReturn() {
	f.SetDefaultHook(func(*logging.Prefix, string, ...interface{}) {
		return
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *LoggerErrorFunc) PushReturn() {
	f.PushHook(func(*logging.Prefix, string, ...interface{}) {
		return
	})
}

func (f *LoggerErrorFunc) nextHook() func(*logging.Prefix, string, ...interface{}) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *LoggerErrorFunc) appendCall(r0 LoggerErrorFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of LoggerErrorFuncCall objects describing the
// invocations of this function.
func (f *LoggerErrorFunc) History() []LoggerErrorFuncCall {
	f.mutex.Lock()
	history := make([]LoggerErrorFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// LoggerErrorFuncCall is an object that describes an invocation of method
// Error on an instance of MockLogger.
type LoggerErrorFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 *logging.Prefix
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 string
	// Arg2 is a slice containing the values of the variadic arguments
	// passed to this method invocation.
	Arg2 []interface{}
}

// Args returns an interface slice containing the arguments of this
// invocation. The variadic slice argument is flattened in this array such
// that one positional argument and three variadic arguments would result in
// a slice of four, not two.
func (c LoggerErrorFuncCall) Args() []interface{} {
	trailing := []interface{}{}
	for _, val := range c.Arg2 {
		trailing = append(trailing, val)
	}

	return append([]interface{}{c.Arg0, c.Arg1}, trailing...)
}

// Results returns an interface slice containing the results of this
// invocation.
func (c LoggerErrorFuncCall) Results() []interface{} {
	return []interface{}{}
}

// LoggerInfoFunc describes the behavior when the Info method of the parent
// MockLogger instance is invoked.
type LoggerInfoFunc struct {
	defaultHook func(*logging.Prefix, string, ...interface{})
	hooks       []func(*logging.Prefix, string, ...interface{})
	history     []LoggerInfoFuncCall
	mutex       sync.Mutex
}

// Info delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockLogger) Info(v0 *logging.Prefix, v1 string, v2 ...interface{}) {
	m.InfoFunc.nextHook()(v0, v1, v2...)
	m.InfoFunc.appendCall(LoggerInfoFuncCall{v0, v1, v2})
	return
}

// SetDefaultHook sets function that is called when the Info method of the
// parent MockLogger instance is invoked and the hook queue is empty.
func (f *LoggerInfoFunc) SetDefaultHook(hook func(*logging.Prefix, string, ...interface{})) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Info method of the parent MockLogger instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *LoggerInfoFunc) PushHook(hook func(*logging.Prefix, string, ...interface{})) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *LoggerInfoFunc) SetDefaultReturn() {
	f.SetDefaultHook(func(*logging.Prefix, string, ...interface{}) {
		return
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *LoggerInfoFunc) PushReturn() {
	f.PushHook(func(*logging.Prefix, string, ...interface{}) {
		return
	})
}

func (f *LoggerInfoFunc) nextHook() func(*logging.Prefix, string, ...interface{}) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *LoggerInfoFunc) appendCall(r0 LoggerInfoFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of LoggerInfoFuncCall objects describing the
// invocations of this function.
func (f *LoggerInfoFunc) History() []LoggerInfoFuncCall {
	f.mutex.Lock()
	history := make([]LoggerInfoFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// LoggerInfoFuncCall is an object that describes an invocation of method
// Info on an instance of MockLogger.
type LoggerInfoFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 *logging.Prefix
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 string
	// Arg2 is a slice containing the values of the variadic arguments
	// passed to this method invocation.
	Arg2 []interface{}
}

// Args returns an interface slice containing the arguments of this
// invocation. The variadic slice argument is flattened in this array such
// that one positional argument and three variadic arguments would result in
// a slice of four, not two.
func (c LoggerInfoFuncCall) Args() []interface{} {
	trailing := []interface{}{}
	for _, val := range c.Arg2 {
		trailing = append(trailing, val)
	}

	return append([]interface{}{c.Arg0, c.Arg1}, trailing...)
}

// Results returns an interface slice containing the results of this
// invocation.
func (c LoggerInfoFuncCall) Results() []interface{} {
	return []interface{}{}
}

// LoggerWarnFunc describes the behavior when the Warn method of the parent
// MockLogger instance is invoked.
type LoggerWarnFunc struct {
	defaultHook func(*logging.Prefix, string, ...interface{})
	hooks       []func(*logging.Prefix, string, ...interface{})
	history     []LoggerWarnFuncCall
	mutex       sync.Mutex
}

// Warn delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockLogger) Warn(v0 *logging.Prefix, v1 string, v2 ...interface{}) {
	m.WarnFunc.nextHook()(v0, v1, v2...)
	m.WarnFunc.appendCall(LoggerWarnFuncCall{v0, v1, v2})
	return
}

// SetDefaultHook sets function that is called when the Warn method of the
// parent MockLogger instance is invoked and the hook queue is empty.
func (f *LoggerWarnFunc) SetDefaultHook(hook func(*logging.Prefix, string, ...interface{})) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Warn method of the parent MockLogger instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *LoggerWarnFunc) PushHook(hook func(*logging.Prefix, string, ...interface{})) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *LoggerWarnFunc) SetDefaultReturn() {
	f.SetDefaultHook(func(*logging.Prefix, string, ...interface{}) {
		return
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *LoggerWarnFunc) PushReturn() {
	f.PushHook(func(*logging.Prefix, string, ...interface{}) {
		return
	})
}

func (f *LoggerWarnFunc) nextHook() func(*logging.Prefix, string, ...interface{}) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *LoggerWarnFunc) appendCall(r0 LoggerWarnFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of LoggerWarnFuncCall objects describing the
// invocations of this function.
func (f *LoggerWarnFunc) History() []LoggerWarnFuncCall {
	f.mutex.Lock()
	history := make([]LoggerWarnFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// LoggerWarnFuncCall is an object that describes an invocation of method
// Warn on an instance of MockLogger.
type LoggerWarnFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 *logging.Prefix
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 string
	// Arg2 is a slice containing the values of the variadic arguments
	// passed to this method invocation.
	Arg2 []interface{}
}

// Args returns an interface slice containing the arguments of this
// invocation. The variadic slice argument is flattened in this array such
// that one positional argument and three variadic arguments would result in
// a slice of four, not two.
func (c LoggerWarnFuncCall) Args() []interface{} {
	trailing := []interface{}{}
	for _, val := range c.Arg2 {
		trailing = append(trailing, val)
	}

	return append([]interface{}{c.Arg0, c.Arg1}, trailing...)
}

// Results returns an interface slice containing the results of this
// invocation.
func (c LoggerWarnFuncCall) Results() []interface{} {
	return []interface{}{}
}
