// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/ij/logging -i ColorPicker -o mock_color_picker_test.go -f

package logging

import "sync"

type MockColorPicker struct {
	ColorizeFunc func(string) string
	histColorize []ColorPickerColorizeParamSet
	mutex        sync.RWMutex
}
type ColorPickerColorizeParamSet struct {
	Arg0 string
}

func NewMockColorPicker() *MockColorPicker {
	m := &MockColorPicker{}
	m.ColorizeFunc = m.defaultColorizeFunc
	return m
}
func (m *MockColorPicker) Colorize(v0 string) string {
	m.mutex.Lock()
	m.histColorize = append(m.histColorize, ColorPickerColorizeParamSet{v0})
	m.mutex.Unlock()
	return m.ColorizeFunc(v0)
}
func (m *MockColorPicker) ColorizeFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histColorize)
}
func (m *MockColorPicker) ColorizeFuncCallParams() []ColorPickerColorizeParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histColorize
}

func (m *MockColorPicker) defaultColorizeFunc(v0 string) string {
	return ""
}
