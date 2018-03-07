package gomap

import (
	"fmt"
)

// ConcurrentMap represents a thread-safe Map.
type ConcurrentMap interface {
	Map
	ClearAsync(callback func())
	ContainsAsync(key interface{}, callback func(bool))
	DeleteAsync(key interface{}, callback func(bool))
	GetAsync(key interface{}, callback func(interface{}, bool))
	LengthAsync(callback func(int))
	SetAsync(key interface{}, value interface{}, callback func(interface{}, bool))
}

// This is a thin wrapper over a thread-safe Map in order to provide additional
// funtionalities (such as async operations).
type concurrentMap struct {
	Map
}

func (cm *concurrentMap) String() string {
	return fmt.Sprint(cm.Map)
}

func (cm *concurrentMap) ContainsAsync(key interface{}, callback func(bool)) {
	go func() {
		found := cm.Contains(key)
		callback(found)
	}()
}

func (cm *concurrentMap) ClearAsync(callback func()) {
	go func() {
		cm.Clear()
		callback()
	}()
}

func (cm *concurrentMap) DeleteAsync(key interface{}, callback func(bool)) {
	go func() {
		result := cm.Delete(key)
		callback(result)
	}()
}

func (cm *concurrentMap) GetAsync(key interface{}, callback func(interface{}, bool)) {
	go func() {
		v, found := cm.Get(key)
		callback(v, found)
	}()
}

func (cm *concurrentMap) LengthAsync(callback func(int)) {
	go func() {
		length := cm.Length()
		callback(length)
	}()
}

func (cm *concurrentMap) SetAsync(key interface{}, value interface{}, callback func(interface{}, bool)) {
	go func() {
		prev, found := cm.Set(key, value)
		callback(prev, found)
	}()
}

func newConcurrentMap(cm Map) ConcurrentMap {
	return &concurrentMap{Map: cm}
}
