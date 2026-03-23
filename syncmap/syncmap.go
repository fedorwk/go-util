package syncmap

import "sync"

type Map[K comparable, V any] struct {
	inner map[K]V
	mu    *sync.RWMutex
}

func New[K comparable, V any](size int) Map[K, V] {
	return Map[K, V]{
		inner: make(map[K]V, size),
		mu:    &sync.RWMutex{},
	}
}

func (m *Map[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inner[key] = value
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.inner[key]
	return val, ok
}

func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.inner, key)
}

func (m *Map[K, V]) Load(replace map[K]V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inner = replace
}

func (m *Map[K, V]) RBorrow() (map[K]V, func()) {
	m.mu.RLock()
	release := func() {
		m.mu.RUnlock()
	}
	return m.inner, release
}
