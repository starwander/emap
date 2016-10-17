// Copyright(c) 2016 Ethan Zhuang <zhuangwj@gmail.com>.

// Package emap implements an enhanced map in golang.
// The main enhancement of emap over original golang map is the support of searching index.
// Values in the emap can have one or more indices which can be used to search or delete.
// Key in the emap must be unique as same as the key in the golang map.
// Index in the emap is an N to M relation which mean a value can have multi indices and multi values can have one same index.
package emap

import (
	"errors"
)

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
