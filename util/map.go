package util

import (
	"sync"
)

// Map Map
type Map struct {
	sync.RWMutex
	m map[interface{}]interface{}
}

func (m *Map) init() {
	if m.m == nil {
		m.m = make(map[interface{}]interface{})
	}
}

// UnsafeGet UnsafeGet
func (m *Map) UnsafeGet(key interface{}) interface{} {
	if m.m == nil {
		return nil
	}
	return m.m[key]
}

// Get Get
func (m *Map) Get(key interface{}) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeGet(key)
}

// UnsafeSet UnsafeSet
func (m *Map) UnsafeSet(key interface{}, value interface{}) {
	m.init()
	m.m[key] = value
}

// Set Set
func (m *Map) Set(key interface{}, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.UnsafeSet(key, value)
}

// TestAndSet TestAndSet
func (m *Map) TestAndSet(key interface{}, value interface{}) interface{} {
	m.Lock()
	defer m.Unlock()

	m.init()

	if v, ok := m.m[key]; ok {
		return v
	}
	m.m[key] = value
	return nil
}

// UnsafeDel UnsafeDel
func (m *Map) UnsafeDel(key interface{}) {
	m.init()
	delete(m.m, key)
}

// Del Del
func (m *Map) Del(key interface{}) {
	m.Lock()
	defer m.Unlock()
	m.UnsafeDel(key)
}

// UnsafeLen UnsafeLen
func (m *Map) UnsafeLen() int {
	if m.m == nil {
		return 0
	}
	return len(m.m)
}

// Len Len
func (m *Map) Len() int {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeLen()
}

// UnsafeRange UnsafeRange
func (m *Map) UnsafeRange(f func(interface{}, interface{}) bool) {
	if m.m == nil {
		return
	}
	for k, v := range m.m {
		if f(k, v) {
			break
		}
	}
}

// RLockRange RLockRange
func (m *Map) RLockRange(f func(interface{}, interface{}) bool) {
	m.RLock()
	defer m.RUnlock()
	m.UnsafeRange(f)
}

// LockRange LockRange
func (m *Map) LockRange(f func(interface{}, interface{}) bool) {
	m.Lock()
	defer m.Unlock()
	m.UnsafeRange(f)
}

// Reset Reset
func (m *Map) Reset() {
	m.m = make(map[interface{}]interface{})
}
