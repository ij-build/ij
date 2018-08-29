// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/ij/registry -i Login -d mocks -f

package mocks

import "sync"

type MockLogin struct {
	GetServerFunc func() (string, error)
	histGetServer []LoginGetServerParamSet
	LoginFunc     func() error
	histLogin     []LoginLoginParamSet
	mutex         sync.RWMutex
}
type LoginGetServerParamSet struct{}
type LoginLoginParamSet struct{}

func NewMockLogin() *MockLogin {
	m := &MockLogin{}
	m.GetServerFunc = m.defaultGetServerFunc
	m.LoginFunc = m.defaultLoginFunc
	return m
}
func (m *MockLogin) GetServer() (string, error) {
	m.mutex.Lock()
	m.histGetServer = append(m.histGetServer, LoginGetServerParamSet{})
	m.mutex.Unlock()
	return m.GetServerFunc()
}
func (m *MockLogin) GetServerFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histGetServer)
}
func (m *MockLogin) GetServerFuncCallParams() []LoginGetServerParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histGetServer
}

func (m *MockLogin) Login() error {
	m.mutex.Lock()
	m.histLogin = append(m.histLogin, LoginLoginParamSet{})
	m.mutex.Unlock()
	return m.LoginFunc()
}
func (m *MockLogin) LoginFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histLogin)
}
func (m *MockLogin) LoginFuncCallParams() []LoginLoginParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histLogin
}

func (m *MockLogin) defaultGetServerFunc() (string, error) {
	return "", nil
}
func (m *MockLogin) defaultLoginFunc() error {
	return nil
}
