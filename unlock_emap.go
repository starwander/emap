package emap

import (
	"errors"
)

type unlockEMap struct {
	store   map[interface{}]interface{}   // key -> value
	keys    map[interface{}][]interface{} // key -> indices
	indices map[interface{}][]interface{} // index -> keys
}

func (m *unlockEMap) KeyNum() int {
	return len(m.keys)
}

func (m *unlockEMap) KeyNumOfIndex(index interface{}) int {
	if keys, exist := m.indices[index]; exist {
		return len(keys)
	}

	return 0
}

func (m *unlockEMap) IndexNum() int {
	return len(m.indices)
}

func (m *unlockEMap) IndexNumOfKey(key interface{}) int {
	if indices, exist := m.keys[key]; exist {
		return len(indices)
	}

	return 0
}

func (m *unlockEMap) HasKey(key interface{}) bool {
	if _, exist := m.keys[key]; exist {
		return true
	}

	return false
}

func (m *unlockEMap) HasIndex(index interface{}) bool {
	if _, exist := m.indices[index]; exist {
		return true
	}

	return false
}

func (m *unlockEMap) Add(key interface{}, value interface{}, indices ...interface{}) error {
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

func (m *unlockEMap) GetByKey(key interface{}) (interface{}, error) {
	if value, exist := m.store[key]; exist {
		return value, nil
	}

	return nil, errors.New("key not exist")
}

func (m *unlockEMap) GetByIndex(index interface{}) ([]interface{}, error) {
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

func (m *unlockEMap) Remove(key interface{}) error {
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

func (m *unlockEMap) AddIndex(key interface{}, index interface{}) error {
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

func (m *unlockEMap) RemoveIndex(key interface{}, index interface{}) error {
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
