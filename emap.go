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
	instance.Store = make(map[interface{}]interface{})
	instance.Keys = make(map[interface{}][]interface{})
	instance.Indices = make(map[interface{}][]interface{})

	return instance
}

func NewExpirableEMap(interval int) EMap {
	instance := new(genericEMap)
	instance.Store = make(map[interface{}]interface{})
	instance.Keys = make(map[interface{}][]interface{})
	instance.Indices = make(map[interface{}][]interface{})

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
	instance.Store = make(map[interface{}]interface{})
	instance.Keys = make(map[interface{}][]interface{})
	instance.Indices = make(map[interface{}][]interface{})

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
	instance.Store = make(map[interface{}]interface{})
	instance.Keys = make(map[interface{}][]interface{})
	instance.Indices = make(map[interface{}][]interface{})

	return instance
}

func add(emap interface{}, key interface{}, value interface{}, indices ...interface{}) error {
	Object := reflect.ValueOf(emap).Elem()
	Store := Object.FieldByName("Store").Interface().(map[interface{}]interface{})
	Keys := Object.FieldByName("Keys").Interface().(map[interface{}][]interface{})
	Indices := Object.FieldByName("Indices").Interface().(map[interface{}][]interface{})

	if _, exist := Keys[key]; exist {
		return errors.New("key duplicte")
	}

	Keys[key] = indices
	Store[key] = value

	for _, index := range indices {
		if keys, exist := Indices[index]; exist {
			Indices[index] = append(keys, key)
		} else {
			Indices[index] = []interface{}{key}
		}
	}

	return nil
}

func fetchByKey(emap interface{}, key interface{}) (interface{}, error) {
	Object := reflect.ValueOf(emap).Elem()
	Store := Object.FieldByName("Store").Interface().(map[interface{}]interface{})

	if value, exist := Store[key]; exist {
		return value, nil
	}

	return nil, errors.New("key not exist")
}

func fetchByIndex(emap interface{}, index interface{}) ([]interface{}, error) {
	Object := reflect.ValueOf(emap).Elem()
	Store := Object.FieldByName("Store").Interface().(map[interface{}]interface{})
	Indices := Object.FieldByName("Indices").Interface().(map[interface{}][]interface{})

	if keys, exist := Indices[index]; exist {
		i := 0
		values := make([]interface{}, len(keys))
		for _, key := range keys {
			values[i] = Store[key]
			i++
		}
		return values, nil
	}

	return nil, errors.New("index not exist")
}

func deleteByKey(emap interface{}, key interface{}) error {
	Object := reflect.ValueOf(emap).Elem()
	Store := Object.FieldByName("Store").Interface().(map[interface{}]interface{})
	Keys := Object.FieldByName("Keys").Interface().(map[interface{}][]interface{})

	if _, exist := Keys[key]; !exist {
		return errors.New("key not exist")
	}

	for _, index := range Keys[key] {
		removeIndex(emap, key, index)
	}

	delete(Keys, key)
	delete(Store, key)

	return nil
}

func deleteByIndex(emap interface{}, index interface{}) error {
	Object := reflect.ValueOf(emap).Elem()
	Indices := Object.FieldByName("Indices").Interface().(map[interface{}][]interface{})

	if _, exist := Indices[index]; !exist {
		return errors.New("index not exist")
	}

	for _, key := range Indices[index] {
		deleteByKey(emap, key)
	}

	return nil
}

func addIndex(emap interface{}, key interface{}, index interface{}) error {
	Object := reflect.ValueOf(emap).Elem()
	Keys := Object.FieldByName("Keys").Interface().(map[interface{}][]interface{})
	Indices := Object.FieldByName("Indices").Interface().(map[interface{}][]interface{})

	if _, exist := Keys[key]; !exist {
		return errors.New("key not exist")
	}

	for _, each := range Keys[key] {
		if each == index {
			return errors.New("index duplicte")
		}
	}
	Keys[key] = append(Keys[key], index)

	if keys, exist := Indices[index]; exist {
		Indices[index] = append(keys, key)
	} else {
		keys = make([]interface{}, 1)
		keys[0] = key
		Indices[index] = keys
	}

	return nil
}

func removeIndex(emap interface{}, key interface{}, index interface{}) error {
	Object := reflect.ValueOf(emap).Elem()
	Keys := Object.FieldByName("Keys").Interface().(map[interface{}][]interface{})
	Indices := Object.FieldByName("Indices").Interface().(map[interface{}][]interface{})

	if _, exist := Keys[key]; !exist {
		return errors.New("key not exist")
	}

	if _, exist := Indices[index]; !exist {
		return errors.New("index not exist")
	}

	for i, each := range Keys[key] {
		if each == index {
			if i == len(Keys[key])-1 {
				Keys[key] = Keys[key][:i]
				break
			}
			Keys[key] = append(Keys[key][:i], Keys[key][i+1:]...)
		}
	}

	for i, each := range Indices[index] {
		if each == key {
			if i == len(Indices[index])-1 {
				Indices[index] = Indices[index][:i]
				break
			}
			Indices[index] = append(Indices[index][:i], Indices[index][i+1:]...)
		}
	}
	if len(Indices[index]) == 0 {
		delete(Indices, index)
	}

	return nil
}

func transform(emap interface{}, callback func(interface{}, interface{}) (interface{}, error)) (map[interface{}]interface{}, error) {
	Object := reflect.ValueOf(emap).Elem()
	Store := Object.FieldByName("Store").Interface().(map[interface{}]interface{})

	var err error
	targets := make(map[interface{}]interface{}, len(Store))

	for key, value := range Store {
		targets[key], err = callback(key, value)
		if err != nil {
			return nil, err
		}
	}

	return targets, nil
}

func foreach(emap interface{}, callback func(interface{}, interface{})) {
	Object := reflect.ValueOf(emap).Elem()
	Store := Object.FieldByName("Store").Interface().(map[interface{}]interface{})

	for key, value := range Store {
		callback(key, value)
	}
}
