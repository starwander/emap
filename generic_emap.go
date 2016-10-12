package emap

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type genericEMap struct {
	mtx      sync.RWMutex
	interval int
	values   map[interface{}]interface{}   // key -> value
	keys     map[interface{}][]interface{} // key -> indices
	indices  map[interface{}][]interface{} // index -> keys
}

func (m *genericEMap) collect(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	for {

		select {
		case <-ticker.C:
			m.mtx.Lock()
			for key, value := range m.values {
				if value.(ExpirableValue).IsExpired() {
					deleteByKey(m.values, m.keys, m.indices, key)
				}
			}
			m.mtx.Unlock()
		}
	}
}

func (m *genericEMap) KeyNum() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.keys)
}

func (m *genericEMap) KeyNumOfIndex(index interface{}) int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if keys, exist := m.indices[index]; exist {
		return len(keys)
	}

	return 0
}

func (m *genericEMap) IndexNum() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.indices)
}

func (m *genericEMap) IndexNumOfKey(key interface{}) int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if indices, exist := m.keys[key]; exist {
		return len(indices)
	}

	return 0
}

func (m *genericEMap) HasKey(key interface{}) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if _, exist := m.keys[key]; exist {
		return true
	}

	return false
}

func (m *genericEMap) HasIndex(index interface{}) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if _, exist := m.indices[index]; exist {
		return true
	}

	return false
}

func (m *genericEMap) Insert(key interface{}, value interface{}, indices ...interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.interval > 0 {
		if _, has := reflect.TypeOf(value).MethodByName("IsExpired"); !has {
			return errors.New("value type wrong")
		}
	}

	return insert(m.values, m.keys, m.indices, key, value, indices...)
}

func (m *genericEMap) FetchByKey(key interface{}) (interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return fetchByKey(m.values, key)
}

func (m *genericEMap) FetchByIndex(index interface{}) ([]interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return fetchByIndex(m.values, m.indices, index)
}

func (m *genericEMap) DeleteByKey(key interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return deleteByKey(m.values, m.keys, m.indices, key)
}

func (m *genericEMap) DeleteByIndex(index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return deleteByIndex(m.values, m.keys, m.indices, index)
}

func (m *genericEMap) AddIndex(key interface{}, index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return addIndex(m.keys, m.indices, key, index)
}

func (m *genericEMap) RemoveIndex(key interface{}, index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return removeIndex(m.keys, m.indices, key, index)
}

func (m *genericEMap) Transform(callback func(interface{}, interface{}) (interface{}, error)) (map[interface{}]interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return transform(m.values, callback)
}

func (m *genericEMap) Foreach(callback func(interface{}, interface{})) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	foreach(m.values, callback)
}
