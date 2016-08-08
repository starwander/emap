package emap

type unlockEMap struct {
	Store   map[interface{}]interface{}   // key -> value
	Keys    map[interface{}][]interface{} // key -> indices
	Indices map[interface{}][]interface{} // index -> keys
}

func (m *unlockEMap) KeyNum() int {
	return len(m.Keys)
}

func (m *unlockEMap) KeyNumOfIndex(index interface{}) int {
	if keys, exist := m.Indices[index]; exist {
		return len(keys)
	}

	return 0
}

func (m *unlockEMap) IndexNum() int {
	return len(m.Indices)
}

func (m *unlockEMap) IndexNumOfKey(key interface{}) int {
	if indices, exist := m.Keys[key]; exist {
		return len(indices)
	}

	return 0
}

func (m *unlockEMap) HasKey(key interface{}) bool {
	if _, exist := m.Keys[key]; exist {
		return true
	}

	return false
}

func (m *unlockEMap) HasIndex(index interface{}) bool {
	if _, exist := m.Indices[index]; exist {
		return true
	}

	return false
}

func (m *unlockEMap) Insert(key interface{}, value interface{}, indices ...interface{}) error {
	return add(m, key, value, indices...)
}

func (m *unlockEMap) FetchByKey(key interface{}) (interface{}, error) {
	return fetchByKey(m, key)
}

func (m *unlockEMap) FetchByIndex(index interface{}) ([]interface{}, error) {
	return fetchByIndex(m, index)
}

func (m *unlockEMap) DeleteByKey(key interface{}) error {
	return deleteByKey(m, key)
}

func (m *unlockEMap) DeleteByIndex(index interface{}) error {
	return deleteByIndex(m, index)
}

func (m *unlockEMap) AddIndex(key interface{}, index interface{}) error {
	return addIndex(m, key, index)
}

func (m *unlockEMap) RemoveIndex(key interface{}, index interface{}) error {
	return removeIndex(m, key, index)
}

func (m *unlockEMap) Transform(callback func(interface{}, interface{}) (interface{}, error)) (map[interface{}]interface{}, error) {
	return transform(m, callback)
}

func (m *unlockEMap) Foreach(callback func(interface{}, interface{}) error) error {
	return foreach(m, callback)
}
