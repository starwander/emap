package emap

import (
	"errors"
	"reflect"
)

type EMap interface {
	Insert(key interface{}, value interface{}, indices ...interface{}) error
	FetchByKey(key interface{}) (interface{}, error)
	FetchByIndex(index interface{}) ([]interface{}, error)
	DeleteByKey(key interface{}) error
	DeleteByIndex(index interface{}) error
	AddIndex(key interface{}, index interface{}) error
	RemoveIndex(key interface{}, index interface{}) error
	KeyNum() int
	KeyNumOfIndex(index interface{}) int
	IndexNum() int
	IndexNumOfKey(key interface{}) int
	HasKey(key interface{}) bool
	HasIndex(index interface{}) bool
	Transform(callback func(interface{}, interface{}) (interface{}, error)) (map[interface{}]interface{}, error)
	Foreach(callback func(interface{}, interface{}))
}

type ExpirableValue interface {
	IsExpired() bool
}

func NewGenericEMap() EMap {
	instance := new(genericEMap)
	instance.values = make(map[interface{}]interface{})
	instance.keys = make(map[interface{}][]interface{})
	instance.indices = make(map[interface{}][]interface{})

	return instance
}

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
