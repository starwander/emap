// Copyright(c) 2016 Ethan Zhuang <zhuangwj@gmail.com>.

// Package emap implements an enhanced map in golang.
// The main enhancement of emap over original golang map is the support of searching index.
// Values in the emap can have one or more indices which can be used to search or delete.
// Key in the emap must be unique as same as the key in the golang map.
// Index in the emap is an N to M relation which mean a value can have multi indices and multi values can have one same index.
package emap

import (
	"errors"
	"reflect"
	"sync"
)

// StrictEMap has a read-write locker inside so it is concurrent safe.
// The type of key, value and index is determined by the input keySample, valueSample and indexSample.
// Only the types of the sample inputs matter, the values of the sample inputs are irrelevant.
// All methods of the strict emap must use the same type of the sample inputs otherwise an error will be returned.
type StrictEMap struct {
	mtx     sync.RWMutex
	values  map[interface{}]interface{}   // key -> value
	keys    map[interface{}][]interface{} // key -> indices
	indices map[interface{}][]interface{} // index -> keys

	keyType     reflect.Kind
	indexType   reflect.Kind
	valueType   reflect.Kind
	valueStruct string
}

// NewStrictEMap creates a new strict emap.
// The types of value, key and index are determined by the inputs.
// Try to appoint any unsupported key or index types, such as pointer, will cause an error return.
func NewStrictEMap(keySample interface{}, valueSample interface{}, indexSample interface{}) (*StrictEMap, error) {
	keyType := reflect.TypeOf(keySample).Kind()
	indexType := reflect.TypeOf(indexSample).Kind()
	valueType := reflect.TypeOf(valueSample).Kind()
	if !isTypeSupported(keyType) || !isTypeSupported(indexType) {
		return nil, errors.New("key or index type not supported")
	}

	instance := new(StrictEMap)
	instance.values = make(map[interface{}]interface{})
	instance.keys = make(map[interface{}][]interface{})
	instance.indices = make(map[interface{}][]interface{})

	instance.keyType = keyType
	instance.indexType = indexType
	instance.valueType = valueType
	if valueType == reflect.Struct {
		instance.valueStruct = reflect.ValueOf(valueSample).Type().Name()
	}

	return instance, nil
}

func isTypeSupported(kind reflect.Kind) bool {
	//if kind == reflect.Int ||
	//kind == reflect.Int8 ||
	//kind == reflect.Int16 ||
	//kind == reflect.Int32 ||
	//kind == reflect.Int64 ||
	//kind == reflect.Uint ||
	//kind == reflect.Uint8 ||
	//kind == reflect.Uint16 ||
	//kind == reflect.Uint32 ||
	//kind == reflect.Uint64 ||
	//kind == reflect.Float32 ||
	//kind == reflect.Float64 ||
	//kind == reflect.Complex64 ||
	//kind == reflect.Complex128 ||
	//kind == reflect.String {
	//	return true
	//}
	if kind == reflect.Bool ||
		kind == reflect.Uintptr ||
		kind == reflect.Array ||
		kind == reflect.Chan ||
		kind == reflect.Func ||
		kind == reflect.Interface ||
		kind == reflect.Map ||
		kind == reflect.Ptr ||
		kind == reflect.Slice ||
		kind == reflect.Struct ||
		kind == reflect.UnsafePointer {
		return false
	}

	return true
}

// KeyNum returns the total key number in the emap.
func (m *StrictEMap) KeyNum() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.keys)
}

// KeyNumOfIndex returns the total key number of the input index in the emap.
func (m *StrictEMap) KeyNumOfIndex(index interface{}) int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if keys, exist := m.indices[index]; exist {
		return len(keys)
	}

	return 0
}

// IndexNum returns the total index number in the emap.
func (m *StrictEMap) IndexNum() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.indices)
}

// IndexNumOfKey returns the total index number of the input key in the emap.
func (m *StrictEMap) IndexNumOfKey(key interface{}) int {
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

// HasKey returns if the input key exists in the emap.
func (m *StrictEMap) HasKey(key interface{}) bool {
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

// HasIndex returns if the input index exists in the emap.
func (m *StrictEMap) HasIndex(index interface{}) bool {
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

// Insert pushes a new value into emap with input key and indices.
// Input key must not be duplicated.
// Input indices are optional.
func (m *StrictEMap) Insert(key interface{}, value interface{}, indices ...interface{}) error {
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
	}
	if m.valueType == reflect.Struct && m.valueStruct != reflect.ValueOf(value).Type().Name() {
		return errors.New("struct type wrong")
	}

	return insert(m.values, m.keys, m.indices, key, value, indices...)
}

// FetchByKey gets the value in the emap by input key.
// Try to fetch a non-existed key will cause an error return.
func (m *StrictEMap) FetchByKey(key interface{}) (interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if m.keyType != reflect.TypeOf(key).Kind() {
		return nil, errors.New("key type wrong")
	}

	return fetchByKey(m.values, key)
}

// FetchByIndex gets the all values in the emap by input index.
// Try to fetch a non-existed index will cause an error return.
func (m *StrictEMap) FetchByIndex(index interface{}) ([]interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if m.indexType != reflect.TypeOf(index).Kind() {
		return nil, errors.New("index type wrong")
	}

	return fetchByIndex(m.values, m.indices, index)
}

// DeleteByKey deletes the value in the emap by input key.
// Try to delete a non-existed key will cause an error return.
func (m *StrictEMap) DeleteByKey(key interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return deleteByKey(m.values, m.keys, m.indices, key)
}

// DeleteByIndex deletes all the values in the emap by input index.
// Try to delete a non-existed index will cause an error return.
func (m *StrictEMap) DeleteByIndex(index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.indexType != reflect.TypeOf(index).Kind() {
		return errors.New("index type wrong")
	}

	return deleteByIndex(m.values, m.keys, m.indices, index)
}

// AddIndex add the input index to the value in the emap of the input key.
// Try to add a duplicate index will cause an error return.
// Try to add an index to a non-existed value will cause an error return.
func (m *StrictEMap) AddIndex(key interface{}, index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.keyType != reflect.TypeOf(key).Kind() {
		return errors.New("key type wrong")
	}

	if m.indexType != reflect.TypeOf(index).Kind() {
		return errors.New("index type wrong")
	}

	return addIndex(m.keys, m.indices, key, index)
}

// RemoveIndex remove the input index from the value in the emap of the input key.
// Try to delete a non-existed index will cause an error return.
// Try to delete an index from a non-existed value will cause an error return.
func (m *StrictEMap) RemoveIndex(key interface{}, index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.keyType != reflect.TypeOf(key).Kind() {
		return errors.New("key type wrong")
	}

	if m.indexType != reflect.TypeOf(index).Kind() {
		return errors.New("index type wrong")
	}

	return removeIndex(m.keys, m.indices, key, index)
}

// Transform is a higher-order operation which apply the input callback function to each key-value pair in the emap.
// Any error returned by the callback function will interrupt the transforming and the error will be returned.
// If transform successfully, a new golang map is created with each key-value pair returned by the input callback function.
func (m *StrictEMap) Transform(callback func(interface{}, interface{}) (interface{}, error)) (map[interface{}]interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return transform(m.values, callback)
}

// Foreach is a higher-order operation which apply the input callback function to each key-value pair in the emap.
// Since the callback function has no return, the foreach procedure will never be interrupted.
// A typical usage of Foreach is apply a closure.
func (m *StrictEMap) Foreach(callback func(interface{}, interface{})) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	foreach(m.values, callback)
}
