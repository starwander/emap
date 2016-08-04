package emap

import (
	"errors"
	"reflect"
	"sync"
)

type strictEMap struct {
	mtx     sync.RWMutex
	Store   map[interface{}]interface{}   // key -> value
	Keys    map[interface{}][]interface{} // key -> indices
	Indices map[interface{}][]interface{} // index -> keys

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

	return len(m.Keys)
}

func (m *strictEMap) KeyNumOfIndex(index interface{}) int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if keys, exist := m.Indices[index]; exist {
		return len(keys)
	}

	return 0
}

func (m *strictEMap) IndexNum() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.Indices)
}

func (m *strictEMap) IndexNumOfKey(key interface{}) int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if m.keyType != reflect.TypeOf(key).Kind() {
		return 0
	}

	if indices, exist := m.Keys[key]; exist {
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

	if _, exist := m.Keys[key]; exist {
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

	if _, exist := m.Indices[index]; exist {
		return true
	}

	return false
}

func (m *strictEMap) Insert(key interface{}, value interface{}, indices ...interface{}) error {
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
		if m.valueType == reflect.Struct && m.valueStruct != reflect.ValueOf(value).Type().Name() {
			return errors.New("struct type wrong")
		}
	}

	return add(m, key, value, indices...)
}

func (m *strictEMap) FetchByKey(key interface{}) (interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if m.keyType != reflect.TypeOf(key).Kind() {
		return nil, errors.New("key type wrong")
	}

	return fetchByKey(m, key)
}

func (m *strictEMap) FetchByIndex(index interface{}) ([]interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if m.indexType != reflect.TypeOf(index).Kind() {
		return nil, errors.New("index type wrong")
	}

	return fetchByIndex(m, index)
}

func (m *strictEMap) DeleteByKey(key interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return deleteByKey(m, key)
}

func (m *strictEMap) DeleteByIndex(index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.indexType != reflect.TypeOf(index).Kind() {
		return errors.New("index type wrong")
	}

	return deleteByIndex(m, index)
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

	return addIndex(m, key, index)
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

	return removeIndex(m, key, index)
}

func (m *strictEMap) Transform(callback func(interface{}, interface{}) (interface{}, error)) (map[interface{}]interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return transform(m, callback)
}
