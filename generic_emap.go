package emap

import (
	"errors"
	"sync"
)

type genericEMap struct {
	mtx     sync.RWMutex
	store   map[interface{}]interface{}   // key -> value
	keys    map[interface{}][]interface{} // key -> indices
	indices map[interface{}][]interface{} // index -> keys
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

func (m *genericEMap) Add(key interface{}, value interface{}, indices ...interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, exist := m.keys[key]; exist {
		return errors.New("key duplicte")
	}

	m.keys[key] = indices
	m.store[key] = value

	for _, index := range indices {
		if keys, exist := m.indices[index]; exist {
			m.indices[index] = append(keys, key)
		} else {
			keys = make([]interface{}, 1)
			keys[0] = key
			m.indices[index] = keys
		}
	}

	return nil
}

func (m *genericEMap) GetByKey(key interface{}) (interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if value, exist := m.store[key]; exist {
		return value, nil
	}

	return nil, errors.New("key not exist")
}

func (m *genericEMap) GetByIndex(index interface{}) ([]interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if keys, exist := m.indices[index]; exist {
		i := 0
		values := make([]interface{}, len(keys))
		for _, key := range keys {
			values[i] = m.store[key]
			i++
		}
		return values, nil
	}

	return nil, errors.New("index not exist")
}

func (m *genericEMap) Remove(key interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, exist := m.keys[key]; !exist {
		return errors.New("key not exist")
	}

	for _, index := range m.keys[key] {
		for i, each := range m.indices[index] {
			if each == key {
				if i == len(m.indices[index])-1 {
					m.indices[index] = m.indices[index][:i]
					break
				}
				m.indices[index] = append(m.indices[index][:i], m.indices[index][i+1:]...)
			}
		}
		if len(m.indices[index]) == 0 {
			delete(m.indices, index)
		}
	}

	delete(m.keys, key)
	delete(m.store, key)

	return nil
}

func (m *genericEMap) AddIndex(key interface{}, index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, exist := m.keys[key]; !exist {
		return errors.New("key not exist")
	}

	for _, each := range m.keys[key] {
		if each == index {
			return errors.New("index duplicte")
		}
	}
	m.keys[key] = append(m.keys[key], index)

	if keys, exist := m.indices[index]; exist {
		m.indices[index] = append(keys, key)
	} else {
		keys = make([]interface{}, 1)
		keys[0] = key
		m.indices[index] = keys
	}

	return nil
}

func (m *genericEMap) RemoveIndex(key interface{}, index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, exist := m.keys[key]; !exist {
		return errors.New("key not exist")
	}

	if _, exist := m.indices[index]; !exist {
		return errors.New("index not exist")
	}

	for i, each := range m.keys[key] {
		if each == index {
			if i == len(m.keys[key])-1 {
				m.keys[key] = m.keys[key][:i]
				break
			}
			m.keys[key] = append(m.keys[key][:i], m.keys[key][i+1:]...)
		}
	}

	for i, each := range m.indices[index] {
		if each == key {
			if i == len(m.indices[index])-1 {
				m.indices[index] = m.indices[index][:i]
				break
			}
			m.indices[index] = append(m.indices[index][:i], m.indices[index][i+1:]...)
		}
	}
	if len(m.indices[index]) == 0 {
		delete(m.indices, index)
	}

	return nil
}
