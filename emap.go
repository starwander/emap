// Copyright(c) 2016 Ethan Zhuang <zhuangwj@gmail.com>.

// Package emap implements an enhanced map in golang.
// The main enhancement of emap over original golang map is the support of searching index.
package emap

import (
	"errors"
	"reflect"
)

// EMap is an enhanced implementation based on the original golang map.
// Values in the emap can have one or more indices which can be used to search or delete.
// Key in the emap must be unique as same as the key in the golang map.
// Index in the emap is an N to M relation which mean a value can have multi indices and multi values can have one same index.
type EMap interface {
	// Insert pushes a new value into emap with input key and indices.
	// Input key must not be duplicated.
	// Input indices are optional.
	Insert(key interface{}, value interface{}, indices ...interface{}) error

	// FetchByKey gets the value in the emap by input key.
	// Try to fetch a non-existed key will cause an error return.
	FetchByKey(key interface{}) (interface{}, error)

	// FetchByIndex gets the all values in the emap by input index.
	// Try to fetch a non-existed index will cause an error return.
	FetchByIndex(index interface{}) ([]interface{}, error)

	// DeleteByKey deletes the value in the emap by input key.
	// Try to delete a non-existed key will cause an error return.
	DeleteByKey(key interface{}) error

	// DeleteByIndex deletes all the values in the emap by input index.
	// Try to delete a non-existed index will cause an error return.
	DeleteByIndex(index interface{}) error

	// AddIndex add the input index to the value in the emap of the input key.
	// Try to add a duplicate index will cause an error return.
	// Try to add an index to a non-existed value will cause an error return.
	AddIndex(key interface{}, index interface{}) error

	// RemoveIndex remove the input index from the value in the emap of the input key.
	// Try to delete a non-existed index will cause an error return.
	// Try to delete an index from a non-existed value will cause an error return.
	RemoveIndex(key interface{}, index interface{}) error

	// KeyNum returns the total key number in the emap.
	KeyNum() int

	// KeyNumOfIndex returns the total key number of the input index in the emap.
	KeyNumOfIndex(index interface{}) int

	// IndexNum returns the total index number in the emap.
	IndexNum() int

	// IndexNumOfKey returns the total index number of the input key in the emap.
	IndexNumOfKey(key interface{}) int

	// HasKey returns if the input key exists in the emap.
	HasKey(key interface{}) bool

	// HasIndex returns if the input index exists in the emap.
	HasIndex(index interface{}) bool

	// Transform is a higher-order operation which apply the input callback function to each key-value pair in the emap.
	// Any error returned by the callback function will interrupt the transforming and the error will be returned.
	// If transform successfully, a new golang map is created with each key-value pair returned by the input callback function.
	Transform(callback func(interface{}, interface{}) (interface{}, error)) (map[interface{}]interface{}, error)

	// Foreach is a higher-order operation which apply the input callback function to each key-value pair in the emap.
	// Since the callback function has no return, the foreach procedure will never be interrupted.
	// A typical usage of Foreach is apply a closure.
	Foreach(callback func(interface{}, interface{}))
}

// ExpirableValue is the interface which must be implemented by all the value in the expirable EMap.
type ExpirableValue interface {
	// IsExpired returns if the value is expired.
	// If true, the value will be deleted automatically.
	IsExpired() bool
}

// NewGenericEMap creates a new emap.
// The generic emap has a read-write locker inside so it is concurrent safe.
func NewGenericEMap() EMap {
	instance := new(genericEMap)
	instance.values = make(map[interface{}]interface{})
	instance.keys = make(map[interface{}][]interface{})
	instance.indices = make(map[interface{}][]interface{})

	return instance
}

// NewExpirableEMap creates a new emap with an expiration checker.
// The expirable emap has a read-write locker inside so it is concurrent safe.
// The expiration checker will check all the values in the emap with the period of input interval(milliseconds).
// All value inserted into the expirable emap must implements ExpirableValue interface of this package.
// If a value is expired, it will be deleted automatically.
func NewExpirableEMap(interval int) EMap {
	instance := new(genericEMap)
	instance.values = make(map[interface{}]interface{})
	instance.keys = make(map[interface{}][]interface{})
	instance.indices = make(map[interface{}][]interface{})

	if interval > 0 {
		instance.interval = interval
		go instance.collect(interval)
	}

	return instance
}

// NewStrictEMap creates a new emap with a type checker.
// The strict emap has a read-write locker inside so it is concurrent safe.
// The type of key, value and index is determined by the input keySample, valueSample and indexSample.
// Only the types of the sample inputs matter, the values of the sample inputs are irrelevant.
// All methods of the strict emap must use the same type of the sample inputs otherwise an error will be returned.
func NewStrictEMap(keySample interface{}, valueSample interface{}, indexSample interface{}) (EMap, error) {
	keyType := reflect.TypeOf(keySample).Kind()
	indexType := reflect.TypeOf(indexSample).Kind()
	valueType := reflect.TypeOf(valueSample).Kind()
	if !isTypeSupported(keyType) || !isTypeSupported(indexType) {
		return nil, errors.New("key or index type not supported")
	}

	instance := new(strictEMap)
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

// NewUnlockEMap creates a new emap without any locker or mutex.
// Since unlock emap is not concurrent safe, it is only suitable for those models like Event Loop to achieve better performance.
func NewUnlockEMap() EMap {
	instance := new(unlockEMap)
	instance.values = make(map[interface{}]interface{})
	instance.keys = make(map[interface{}][]interface{})
	instance.indices = make(map[interface{}][]interface{})

	return instance
}

func insert(valueStore map[interface{}]interface{}, keyStore map[interface{}][]interface{}, indexStore map[interface{}][]interface{}, key interface{}, value interface{}, indices ...interface{}) error {
	if _, exist := keyStore[key]; exist {
		return errors.New("key duplicte")
	}

	keyStore[key] = indices
	valueStore[key] = value

	for _, index := range indices {
		if keys, exist := indexStore[index]; exist {
			indexStore[index] = append(keys, key)
		} else {
			indexStore[index] = []interface{}{key}
		}
	}

	return nil
}

func fetchByKey(valueStore map[interface{}]interface{}, key interface{}) (interface{}, error) {
	if value, exist := valueStore[key]; exist {
		return value, nil
	}

	return nil, errors.New("key not exist")
}

func fetchByIndex(valueStore map[interface{}]interface{}, indexStore map[interface{}][]interface{}, index interface{}) ([]interface{}, error) {
	if keys, exist := indexStore[index]; exist {
		values := make([]interface{}, len(keys))
		for i, key := range keys {
			values[i] = valueStore[key]
		}
		return values, nil
	}

	return nil, errors.New("index not exist")
}

func deleteByKey(valueStore map[interface{}]interface{}, keyStore map[interface{}][]interface{}, indexStore map[interface{}][]interface{}, key interface{}) error {
	if _, exist := keyStore[key]; !exist {
		return errors.New("key not exist")
	}

	for _, index := range keyStore[key] {
		removeIndex(keyStore, indexStore, key, index)
	}

	delete(keyStore, key)
	delete(valueStore, key)

	return nil
}

func deleteByIndex(valueStore map[interface{}]interface{}, keyStore map[interface{}][]interface{}, indexStore map[interface{}][]interface{}, index interface{}) error {
	if _, exist := indexStore[index]; !exist {
		return errors.New("index not exist")
	}

	for _, key := range indexStore[index] {
		deleteByKey(valueStore, keyStore, indexStore, key)
	}

	return nil
}

func addIndex(keyStore map[interface{}][]interface{}, indexStore map[interface{}][]interface{}, key interface{}, index interface{}) error {
	if _, exist := keyStore[key]; !exist {
		return errors.New("key not exist")
	}

	for _, each := range keyStore[key] {
		if each == index {
			return errors.New("index duplicte")
		}
	}
	keyStore[key] = append(keyStore[key], index)

	if keys, exist := indexStore[index]; exist {
		indexStore[index] = append(keys, key)
	} else {
		indexStore[index] = []interface{}{key}
	}

	return nil
}

func removeIndex(keyStore map[interface{}][]interface{}, indexStore map[interface{}][]interface{}, key interface{}, index interface{}) error {
	if _, exist := keyStore[key]; !exist {
		return errors.New("key not exist")
	}

	if _, exist := indexStore[index]; !exist {
		return errors.New("index not exist")
	}

	for i, each := range keyStore[key] {
		if each == index {
			if i == len(keyStore[key])-1 {
				keyStore[key] = keyStore[key][:i]
			} else {
				keyStore[key] = append(keyStore[key][:i], keyStore[key][i+1:]...)
			}
			break
		}
	}

	for i, each := range indexStore[index] {
		if each == key {
			if i == len(indexStore[index])-1 {
				indexStore[index] = indexStore[index][:i]
			} else {
				indexStore[index] = append(indexStore[index][:i], indexStore[index][i+1:]...)
			}
			break
		}
	}
	if len(indexStore[index]) == 0 {
		delete(indexStore, index)
	}

	return nil
}

func transform(valueStore map[interface{}]interface{}, callback func(interface{}, interface{}) (interface{}, error)) (map[interface{}]interface{}, error) {
	var err error
	targets := make(map[interface{}]interface{}, len(valueStore))

	for key, value := range valueStore {
		targets[key], err = callback(key, value)
		if err != nil {
			return nil, err
		}
	}

	return targets, nil
}

func foreach(valueStore map[interface{}]interface{}, callback func(interface{}, interface{})) {
	for key, value := range valueStore {
		callback(key, value)
	}
}
