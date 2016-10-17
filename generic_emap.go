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

// GenericEMap has a read-write locker inside so it is concurrent safe.
// The value, key and index type is unlimited in the generic emap.
type GenericEMap struct {
	mtx      sync.RWMutex
	interval int
	values   map[interface{}]interface{}   // key -> value
	keys     map[interface{}][]interface{} // key -> indices
	indices  map[interface{}][]interface{} // index -> keys
}

// NewGenericEMap creates a new generic emap.
func NewGenericEMap() *GenericEMap {
	instance := new(GenericEMap)
	instance.values = make(map[interface{}]interface{})
	instance.keys = make(map[interface{}][]interface{})
	instance.indices = make(map[interface{}][]interface{})

	return instance
}

// KeyNum returns the total key number in the emap.
func (m *GenericEMap) KeyNum() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.keys)
}

// KeyNumOfIndex returns the total key number of the input index in the emap.
func (m *GenericEMap) KeyNumOfIndex(index interface{}) int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if keys, exist := m.indices[index]; exist {
		return len(keys)
	}

	return 0
}

// IndexNum returns the total index number in the emap.
func (m *GenericEMap) IndexNum() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.indices)
}

// IndexNumOfKey returns the total index number of the input key in the emap.
func (m *GenericEMap) IndexNumOfKey(key interface{}) int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if indices, exist := m.keys[key]; exist {
		return len(indices)
	}

	return 0
}

// HasKey returns if the input key exists in the emap.
func (m *GenericEMap) HasKey(key interface{}) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if _, exist := m.keys[key]; exist {
		return true
	}

	return false
}

// HasIndex returns if the input index exists in the emap.
func (m *GenericEMap) HasIndex(index interface{}) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if _, exist := m.indices[index]; exist {
		return true
	}

	return false
}

// Insert pushes a new value into emap with input key and indices.
// Input key must not be duplicated.
// Input indices are optional.
func (m *GenericEMap) Insert(key interface{}, value interface{}, indices ...interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.interval > 0 {
		if _, has := reflect.TypeOf(value).MethodByName("IsExpired"); !has {
			return errors.New("value type wrong")
		}
	}

	return insert(m.values, m.keys, m.indices, key, value, indices...)
}

// FetchByKey gets the value in the emap by input key.
// Try to fetch a non-existed key will cause an error return.
func (m *GenericEMap) FetchByKey(key interface{}) (interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return fetchByKey(m.values, key)
}

// FetchByIndex gets the all values in the emap by input index.
// Try to fetch a non-existed index will cause an error return.
func (m *GenericEMap) FetchByIndex(index interface{}) ([]interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return fetchByIndex(m.values, m.indices, index)
}

// DeleteByKey deletes the value in the emap by input key.
// Try to delete a non-existed key will cause an error return.
func (m *GenericEMap) DeleteByKey(key interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return deleteByKey(m.values, m.keys, m.indices, key)
}

// DeleteByIndex deletes all the values in the emap by input index.
// Try to delete a non-existed index will cause an error return.
func (m *GenericEMap) DeleteByIndex(index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return deleteByIndex(m.values, m.keys, m.indices, index)
}

// AddIndex add the input index to the value in the emap of the input key.
// Try to add a duplicate index will cause an error return.
// Try to add an index to a non-existed value will cause an error return.
func (m *GenericEMap) AddIndex(key interface{}, index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return addIndex(m.keys, m.indices, key, index)
}

// RemoveIndex remove the input index from the value in the emap of the input key.
// Try to delete a non-existed index will cause an error return.
// Try to delete an index from a non-existed value will cause an error return.
func (m *GenericEMap) RemoveIndex(key interface{}, index interface{}) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return removeIndex(m.keys, m.indices, key, index)
}

// Transform is a higher-order operation which apply the input callback function to each key-value pair in the emap.
// Any error returned by the callback function will interrupt the transforming and the error will be returned.
// If transform successfully, a new golang map is created with each key-value pair returned by the input callback function.
func (m *GenericEMap) Transform(callback func(interface{}, interface{}) (interface{}, error)) (map[interface{}]interface{}, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return transform(m.values, callback)
}

// Foreach is a higher-order operation which apply the input callback function to each key-value pair in the emap.
// Since the callback function has no return, the foreach procedure will never be interrupted.
// A typical usage of Foreach is apply a closure.
func (m *GenericEMap) Foreach(callback func(interface{}, interface{})) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	foreach(m.values, callback)
}
