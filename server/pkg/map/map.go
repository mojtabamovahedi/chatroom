package chatMap

import "sync"

type Map[T comparable, V any] struct {
	m    map[T]V
	lock sync.RWMutex
}

func NewMap[T comparable, V any]() *Map[T, V] {
	return &Map[T, V]{
		m: make(map[T]V),
	}
}

func (m *Map[T, V]) Set(key T, value V) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	m.m[key] = value
}

func (m *Map[T, V]) Get(key T) (V, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	value, ok := m.m[key]
	return value, ok
}

func (m *Map[T, V]) Remove(key T) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.m, key)
}
