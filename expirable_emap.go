// Copyright(c) 2016 Ethan Zhuang <zhuangwj@gmail.com>.

package emap

import (
	"time"
)

// ExpirableValue is the interface which must be implemented by all the value in the expirable EMap.
type ExpirableValue interface {
	// IsExpired returns if the value is expired.
	// If true, the value will be deleted automatically.
	IsExpired() bool
}

// NewExpirableEMap creates a new generic emap with an expiration checker.
// The expiration checker will check all the values in the emap with the period of input interval(milliseconds).
// All value inserted into the expirable emap must implements ExpirableValue interface of this package.
// If a value is expired, it will be deleted automatically.
func NewExpirableEMap(interval int) *GenericEMap {
	instance := new(GenericEMap)
	instance.values = make(map[interface{}]interface{})
	instance.keys = make(map[interface{}][]interface{})
	instance.indices = make(map[interface{}][]interface{})

	if interval > 0 {
		instance.interval = interval
		go instance.collect(interval)
	}

	return instance
}

func (m *GenericEMap) collect(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			m.mtx.Lock()
			for key, value := range m.values {
				if value.(ExpirableValue).IsExpired() {
					deleteByKey(m.values, m.keys, m.indices, key)
				}
			}
			m.mtx.Unlock()
		}
	}
}
