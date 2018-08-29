// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/ij/registry -i Login -d mocks -f

package mocks

import "sync"

type MockLogin struct {
	LoginFunc func() (string, error)
	histLogin []LoginLoginParamSet
	mutex     sync.RWMutex
}
type LoginLoginParamSet struct{}

func NewMockLogin() *MockLogin {
	m := &MockLogin{}
	m.LoginFunc = m.defaultLoginFunc
	return m
}
func (m *MockLogin) Login() (string, error) {
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

func (m *MockLogin) defaultLoginFunc() (string, error) {
	return "", nil
}
