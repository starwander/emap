package emap

import (
	"errors"
	"reflect"
)

type EMap interface {
	Add(key interface{}, value interface{}, indices ...interface{}) error
	Remove(key interface{}) error
	GetByKey(key interface{}) (interface{}, error)
	GetByIndex(index interface{}) ([]interface{}, error)
	AddIndex(key interface{}, index interface{}) error
	RemoveIndex(key interface{}, index interface{}) error
	KeyNum() int
	KeyNumOfIndex(index interface{}) int
	IndexNum() int
	IndexNumOfKey(key interface{}) int
	HasKey(key interface{}) bool
	HasIndex(index interface{}) bool
}

func NewGenericEMap() EMap {
	instance := new(genericEMap)
	instance.store = make(map[interface{}]interface{})
	instance.keys = make(map[interface{}][]interface{})
	instance.indices = make(map[interface{}][]interface{})

	return instance
}

func NewExpirableEMap(interval int) EMap {
	instance := new(expirableEMap)
	instance.store = make(map[interface{}]interface{})
	instance.keys = make(map[interface{}][]interface{})
	instance.indices = make(map[interface{}][]interface{})

	go instance.collect(interval)

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
	instance.store = make(map[interface{}]interface{})
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
	instance.store = make(map[interface{}]interface{})
	instance.keys = make(map[interface{}][]interface{})
	instance.indices = make(map[interface{}][]interface{})

	return instance
}
