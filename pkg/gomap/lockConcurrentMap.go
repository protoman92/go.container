package gomap

import (
	"fmt"
	"sync"
)

type lockConcurrentMap struct {
	mutex   *sync.RWMutex
	storage Map
}

func (lcm *lockConcurrentMap) String() string {
	lcm.mutex.RLock()
	defer lcm.mutex.RUnlock()
	return fmt.Sprint(lcm.storage)
}

func (lcm *lockConcurrentMap) Clear() {
	lcm.mutex.Lock()
	defer lcm.mutex.Unlock()
	lcm.storage.Clear()
}

func (lcm *lockConcurrentMap) Contains(key interface{}) bool {
	lcm.mutex.RLock()
	defer lcm.mutex.RUnlock()
	return lcm.storage.Contains(key)
}

func (lcm *lockConcurrentMap) Delete(key interface{}) bool {
	lcm.mutex.Lock()
	defer lcm.mutex.Unlock()
	return lcm.storage.Delete(key)
}

func (lcm *lockConcurrentMap) Get(key interface{}) (interface{}, bool) {
	lcm.mutex.RLock()
	defer lcm.mutex.RUnlock()
	return lcm.storage.Get(key)
}

func (lcm *lockConcurrentMap) Length() int {
	lcm.mutex.RLock()
	defer lcm.mutex.RUnlock()
	return lcm.storage.Length()
}

func (lcm *lockConcurrentMap) Keys() []interface{} {
	lcm.mutex.RLock()
	defer lcm.mutex.RUnlock()
	return lcm.storage.Keys()
}

func (lcm *lockConcurrentMap) Set(key interface{}, value interface{}) (interface{}, bool) {
	lcm.mutex.Lock()
	defer lcm.mutex.Unlock()
	return lcm.storage.Set(key, value)
}

// NewLockConcurrentMap returns a new lock-based ConcurrentMap.
func NewLockConcurrentMap(storage Map) ConcurrentMap {
	cm := &lockConcurrentMap{mutex: &sync.RWMutex{}, storage: storage}
	return newConcurrentMap(cm)
}
