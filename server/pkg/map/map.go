package chatMap

import "sync"

// Map is a thread-safe map
type Map[T comparable, V any] struct {
    m    map[T]V
    lock sync.RWMutex
}

// NewMap creates a new thread-safe map
func NewMap[T comparable, V any]() *Map[T, V] {
    return &Map[T, V]{
        m: make(map[T]V),
    }
}

// Set adds or updates a key-value pair in the map
func (m *Map[T, V]) Set(key T, value V) {
    m.lock.RLock()
    defer m.lock.RUnlock()
    m.m[key] = value
}

// Get retrieves a value by key from the map
func (m *Map[T, V]) Get(key T) (V, bool) {
    m.lock.RLock()
    defer m.lock.RUnlock()
    value, ok := m.m[key]
    return value, ok
}

// Remove deletes a key-value pair from the map
func (m *Map[T, V]) Remove(key T) {
    m.lock.Lock()
    defer m.lock.Unlock()
    delete(m.m, key)
}