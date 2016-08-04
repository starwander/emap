package emap

import (
	"sync"
)

type genericEMap struct {
	mtx     sync.RWMutex
	Store   map[interface{}]interface{}   // key -> value
	Keys    map[interface{}][]interface{} // key -> indices
	Indices map[interface{}][]interface{} // index -> keys
}

func (m *genericEMap) KeyNum() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.Keys)
}

func (m *genericEMap) KeyNumOfIndex(index interface{}) int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if keys, exist := m.Indices[index]; exist {
		return len(keys)
	}

	return 0
}

func (m *genericEMap) IndexNum() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.Indices)
}

func (m *genericEMap) IndexNumOfKey(key interface{}) int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if indices, exist := m.Keys[key]; exist {
		return len(indices)
	}

	return 0
}

func (m *genericEMap) HasKey(key interface{}) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if _, exist := m.Keys[key]; exist {
		return true
	}

	return false
}

func (m *genericEMap) HasIndex(index interface{}) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if _, exist := m.Indices[index]; exist {
		return true
	}

	return false
}

func (m *genericEMap) Insert(key interface{}, value interface{}, indices ...interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return add(m, key, value, indices...)
}

func (m *genericEMap) FetchByKey(key interface{}) (interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return fetchByKey(m, key)
}

func (m *genericEMap) FetchByIndex(index interface{}) ([]interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return fetchByIndex(m, index)
}

func (m *genericEMap) DeleteByKey(key interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return deleteByKey(m, key)
}

func (m *genericEMap) DeleteByIndex(index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return deleteByIndex(m, index)
}

func (m *genericEMap) AddIndex(key interface{}, index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return addIndex(m, key, index)
}

func (m *genericEMap) RemoveIndex(key interface{}, index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return removeIndex(m, key, index)
}

func (m *genericEMap) Transform(callback func(interface{}, interface{}) (interface{}, error)) (map[interface{}]interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return transform(m, callback)
}
