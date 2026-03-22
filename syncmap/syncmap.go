package syncmap

import "sync"

type SyncMap[K comparable, V any] struct {
	inner map[K]V
	mu    *sync.RWMutex
}

func NewSyncMap[K comparable, V any](size int) SyncMap[K, V] {
	return SyncMap[K, V]{
		inner: make(map[K]V, size),
		mu:    &sync.RWMutex{},
	}
}

func (m *SyncMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inner[key] = value
}

func (m *SyncMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.inner[key]
	return val, ok
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.inner, key)
}

func (m *SyncMap[K, V]) Load(replace map[K]V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inner = replace
}

func (m *SyncMap[K, V]) RBorrow() (map[K]V, func()) {
	m.mu.RLock()
	release := func() {
		m.mu.RUnlock()
	}
	return m.inner, release
}
