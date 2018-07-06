// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/gache -i Cache -d mocks

package mocks

import (
	gache "github.com/efritz/gache"
	sync "sync"
)

type MockCache struct {
	statsBustTagsLock          sync.RWMutex
	statBustTagsFuncCallCount  int
	statBustTagsFuncCallParams []CacheBustTagsParamSet
	BustTagsFunc               func(...string) error

	statsGetValueLock          sync.RWMutex
	statGetValueFuncCallCount  int
	statGetValueFuncCallParams []CacheGetValueParamSet
	GetValueFunc               func(string) (string, error)

	statsRemoveLock          sync.RWMutex
	statRemoveFuncCallCount  int
	statRemoveFuncCallParams []CacheRemoveParamSet
	RemoveFunc               func(string) error

	statsSetValueLock          sync.RWMutex
	statSetValueFuncCallCount  int
	statSetValueFuncCallParams []CacheSetValueParamSet
	SetValueFunc               func(string, string, ...string) error
}
type CacheRemoveParamSet struct {
	Arg0 string
}
type CacheSetValueParamSet struct {
	Arg0 string
	Arg1 string
	Arg2 []string
}
type CacheBustTagsParamSet struct {
	Arg0 []string
}
type CacheGetValueParamSet struct {
	Arg0 string
}

var _ gache.Cache = NewMockCache()

func NewMockCache() *MockCache {
	m := &MockCache{}
	m.SetValueFunc = m.defaultSetValueFunc
	m.BustTagsFunc = m.defaultBustTagsFunc
	m.GetValueFunc = m.defaultGetValueFunc
	m.RemoveFunc = m.defaultRemoveFunc
	return m
}
func (m *MockCache) Remove(v0 string) error {
	m.statsRemoveLock.Lock()
	m.statRemoveFuncCallCount++
	m.statRemoveFuncCallParams = append(m.statRemoveFuncCallParams, CacheRemoveParamSet{v0})
	m.statsRemoveLock.Unlock()
	return m.RemoveFunc(v0)
}
func (m *MockCache) RemoveFuncCallCount() int {
	m.statsRemoveLock.RLock()
	defer m.statsRemoveLock.RUnlock()
	return m.statRemoveFuncCallCount
}
func (m *MockCache) RemoveFuncCallParams() []CacheRemoveParamSet {
	m.statsRemoveLock.RLock()
	defer m.statsRemoveLock.RUnlock()
	return m.statRemoveFuncCallParams
}

func (m *MockCache) SetValue(v0 string, v1 string, v2 ...string) error {
	m.statsSetValueLock.Lock()
	m.statSetValueFuncCallCount++
	m.statSetValueFuncCallParams = append(m.statSetValueFuncCallParams, CacheSetValueParamSet{v0, v1, v2})
	m.statsSetValueLock.Unlock()
	return m.SetValueFunc(v0, v1, v2...)
}
func (m *MockCache) SetValueFuncCallCount() int {
	m.statsSetValueLock.RLock()
	defer m.statsSetValueLock.RUnlock()
	return m.statSetValueFuncCallCount
}
func (m *MockCache) SetValueFuncCallParams() []CacheSetValueParamSet {
	m.statsSetValueLock.RLock()
	defer m.statsSetValueLock.RUnlock()
	return m.statSetValueFuncCallParams
}

func (m *MockCache) BustTags(v0 ...string) error {
	m.statsBustTagsLock.Lock()
	m.statBustTagsFuncCallCount++
	m.statBustTagsFuncCallParams = append(m.statBustTagsFuncCallParams, CacheBustTagsParamSet{v0})
	m.statsBustTagsLock.Unlock()
	return m.BustTagsFunc(v0...)
}
func (m *MockCache) BustTagsFuncCallCount() int {
	m.statsBustTagsLock.RLock()
	defer m.statsBustTagsLock.RUnlock()
	return m.statBustTagsFuncCallCount
}
func (m *MockCache) BustTagsFuncCallParams() []CacheBustTagsParamSet {
	m.statsBustTagsLock.RLock()
	defer m.statsBustTagsLock.RUnlock()
	return m.statBustTagsFuncCallParams
}

func (m *MockCache) GetValue(v0 string) (string, error) {
	m.statsGetValueLock.Lock()
	m.statGetValueFuncCallCount++
	m.statGetValueFuncCallParams = append(m.statGetValueFuncCallParams, CacheGetValueParamSet{v0})
	m.statsGetValueLock.Unlock()
	return m.GetValueFunc(v0)
}
func (m *MockCache) GetValueFuncCallCount() int {
	m.statsGetValueLock.RLock()
	defer m.statsGetValueLock.RUnlock()
	return m.statGetValueFuncCallCount
}
func (m *MockCache) GetValueFuncCallParams() []CacheGetValueParamSet {
	m.statsGetValueLock.RLock()
	defer m.statsGetValueLock.RUnlock()
	return m.statGetValueFuncCallParams
}

func (m *MockCache) defaultBustTagsFunc(v0 ...string) error {
	return nil
}
func (m *MockCache) defaultGetValueFunc(v0 string) (string, error) {
	return "", nil
}
func (m *MockCache) defaultRemoveFunc(v0 string) error {
	return nil
}
func (m *MockCache) defaultSetValueFunc(v0 string, v1 string, v2 ...string) error {
	return nil
}
