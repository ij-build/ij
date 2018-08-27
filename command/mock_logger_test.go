// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/ij/logging -i Logger -o mock_logger_test.go -f

package command

import (
	logging "github.com/efritz/ij/logging"
	"sync"
)

type MockLogger struct {
	DebugFunc func(*logging.Prefix, string, ...interface{})
	histDebug []LoggerDebugParamSet
	ErrorFunc func(*logging.Prefix, string, ...interface{})
	histError []LoggerErrorParamSet
	InfoFunc  func(*logging.Prefix, string, ...interface{})
	histInfo  []LoggerInfoParamSet
	mutex     sync.RWMutex
}
type LoggerDebugParamSet struct {
	Arg0 *logging.Prefix
	Arg1 string
	Arg2 []interface{}
}
type LoggerErrorParamSet struct {
	Arg0 *logging.Prefix
	Arg1 string
	Arg2 []interface{}
}
type LoggerInfoParamSet struct {
	Arg0 *logging.Prefix
	Arg1 string
	Arg2 []interface{}
}

func NewMockLogger() *MockLogger {
	m := &MockLogger{}
	m.DebugFunc = m.defaultDebugFunc
	m.ErrorFunc = m.defaultErrorFunc
	m.InfoFunc = m.defaultInfoFunc
	return m
}
func (m *MockLogger) Debug(v0 *logging.Prefix, v1 string, v2 ...interface{}) {
	m.mutex.Lock()
	m.histDebug = append(m.histDebug, LoggerDebugParamSet{v0, v1, v2})
	m.mutex.Unlock()
	m.DebugFunc(v0, v1, v2...)
}
func (m *MockLogger) DebugFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histDebug)
}
func (m *MockLogger) DebugFuncCallParams() []LoggerDebugParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histDebug
}

func (m *MockLogger) Error(v0 *logging.Prefix, v1 string, v2 ...interface{}) {
	m.mutex.Lock()
	m.histError = append(m.histError, LoggerErrorParamSet{v0, v1, v2})
	m.mutex.Unlock()
	m.ErrorFunc(v0, v1, v2...)
}
func (m *MockLogger) ErrorFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histError)
}
func (m *MockLogger) ErrorFuncCallParams() []LoggerErrorParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histError
}

func (m *MockLogger) Info(v0 *logging.Prefix, v1 string, v2 ...interface{}) {
	m.mutex.Lock()
	m.histInfo = append(m.histInfo, LoggerInfoParamSet{v0, v1, v2})
	m.mutex.Unlock()
	m.InfoFunc(v0, v1, v2...)
}
func (m *MockLogger) InfoFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histInfo)
}
func (m *MockLogger) InfoFuncCallParams() []LoggerInfoParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histInfo
}

func (m *MockLogger) defaultDebugFunc(v0 *logging.Prefix, v1 string, v2 ...interface{}) {
	return
}
func (m *MockLogger) defaultErrorFunc(v0 *logging.Prefix, v1 string, v2 ...interface{}) {
	return
}
func (m *MockLogger) defaultInfoFunc(v0 *logging.Prefix, v1 string, v2 ...interface{}) {
	return
}
