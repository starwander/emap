package emap

import (
	"errors"
	"reflect"
	"sync"
)

type strictEMap struct {
	mtx     sync.RWMutex
	store   map[interface{}]interface{}   // key -> value
	keys    map[interface{}][]interface{} // key -> indices
	indices map[interface{}][]interface{} // index -> keys

	keyType     reflect.Kind
	indexType   reflect.Kind
	valueType   reflect.Kind
	valueStruct string
}

func isTypeSupported(kind reflect.Kind) bool {
	if kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64 ||
		kind == reflect.Float32 ||
		kind == reflect.Float64 ||
		kind == reflect.Complex64 ||
		kind == reflect.Complex128 ||
		kind == reflect.Complex128 ||
		kind == reflect.String {
		return true
	}

	return false
}

func (m *strictEMap) KeyNum() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.keys)
}

func (m *strictEMap) KeyNumOfIndex(index interface{}) int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if keys, exist := m.indices[index]; exist {
		return len(keys)
	}

	return 0
}

func (m *strictEMap) IndexNum() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.indices)
}

func (m *strictEMap) IndexNumOfKey(key interface{}) int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if m.keyType != reflect.TypeOf(key).Kind() {
		return 0
	}

	if indices, exist := m.keys[key]; exist {
		return len(indices)
	}

	return 0
}

func (m *strictEMap) HasKey(key interface{}) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if m.keyType != reflect.TypeOf(key).Kind() {
		return false
	}

	if _, exist := m.keys[key]; exist {
		return true
	}

	return false
}

func (m *strictEMap) HasIndex(index interface{}) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if m.indexType != reflect.TypeOf(index).Kind() {
		return false
	}

	if _, exist := m.indices[index]; exist {
		return true
	}

	return false
}

func (m *strictEMap) Add(key interface{}, value interface{}, indices ...interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.keyType != reflect.TypeOf(key).Kind() {
		return errors.New("key type wrong")
	}
	for _, index := range indices {
		if m.indexType != reflect.TypeOf(index).Kind() {
			return errors.New("index type wrong")
		}

	}
	if m.valueType != reflect.TypeOf(value).Kind() {
		return errors.New("value type wrong")
	} else {
		if m.valueStruct != reflect.ValueOf(value).Type().Name() {
			return errors.New("struct type wrong")
		}
	}

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

func (m *strictEMap) GetByKey(key interface{}) (interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if m.keyType != reflect.TypeOf(key).Kind() {
		return nil, errors.New("key type wrong")
	}

	if value, exist := m.store[key]; exist {
		return value, nil
	}

	return nil, errors.New("key not exist")
}

func (m *strictEMap) GetByIndex(index interface{}) ([]interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if m.indexType != reflect.TypeOf(index).Kind() {
		return nil, errors.New("index type wrong")
	}

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

func (m *strictEMap) Remove(key interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return m.remove(key)
}

func (m *strictEMap) remove(key interface{}) error {
	if m.keyType != reflect.TypeOf(key).Kind() {
		return errors.New("key type wrong")
	}

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

func (m *strictEMap) AddIndex(key interface{}, index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.keyType != reflect.TypeOf(key).Kind() {
		return errors.New("key type wrong")
	}

	if m.indexType != reflect.TypeOf(index).Kind() {
		return errors.New("index type wrong")
	}

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

func (m *strictEMap) RemoveIndex(key interface{}, index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.keyType != reflect.TypeOf(key).Kind() {
		return errors.New("key type wrong")
	}

	if m.indexType != reflect.TypeOf(index).Kind() {
		return errors.New("index type wrong")
	}

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
