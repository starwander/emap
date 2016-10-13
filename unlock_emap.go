// Copyright(c) 2016 Ethan Zhuang <zhuangwj@gmail.com>.

package emap

type unlockEMap struct {
	values  map[interface{}]interface{}   // key -> value
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

func (m *unlockEMap) Insert(key interface{}, value interface{}, indices ...interface{}) error {
	return insert(m.values, m.keys, m.indices, key, value, indices...)
}

func (m *unlockEMap) FetchByKey(key interface{}) (interface{}, error) {
	return fetchByKey(m.values, key)
}

func (m *unlockEMap) FetchByIndex(index interface{}) ([]interface{}, error) {
	return fetchByIndex(m.values, m.indices, index)
}

func (m *unlockEMap) DeleteByKey(key interface{}) error {
	return deleteByKey(m.values, m.keys, m.indices, key)
}

func (m *unlockEMap) DeleteByIndex(index interface{}) error {
	return deleteByIndex(m.values, m.keys, m.indices, index)
}

func (m *unlockEMap) AddIndex(key interface{}, index interface{}) error {
	return addIndex(m.keys, m.indices, key, index)
}

func (m *unlockEMap) RemoveIndex(key interface{}, index interface{}) error {
	return removeIndex(m.keys, m.indices, key, index)
}

func (m *unlockEMap) Transform(callback func(interface{}, interface{}) (interface{}, error)) (map[interface{}]interface{}, error) {
	return transform(m.values, callback)
}

func (m *unlockEMap) Foreach(callback func(interface{}, interface{})) {
	foreach(m.values, callback)
}
