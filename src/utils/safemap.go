package utils

import (
	"sync"
)

type SafeMap struct {
	mu      *sync.RWMutex
	currMap map[interface{}]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		mu:      new(sync.RWMutex),
		currMap: make(map[interface{}]interface{}),
	}
}

func (m *SafeMap) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.currMap)
}

func (m *SafeMap) Get(k interface{}) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if v, ok := m.currMap[k]; ok {
		return v
	}
	return nil
}

func (m *SafeMap) Set(k interface{}, v interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currMap[k] = v
}

func (m *SafeMap) Delete(k interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.currMap, k)
}

func (m *SafeMap) Items() map[interface{}]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	newMap := make(map[interface{}]interface{})
	for k, v := range m.currMap {
		newMap[k] = v
	}
	return newMap
}
