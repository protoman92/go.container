package gomap

// ConcurrentMap represents a thread-safe Map.
type ConcurrentMap interface {
	Map
	UnderlyingStorageAsync(callback func(map[Key]Value))
	ClearAsync(callback func())
	ContainsAsync(key Key, callback func(bool))
	DeleteAsync(key Key, callback func(int))
	GetAsync(key Key, callback func(Value, bool))
	IsEmptyAsync(callback func(bool))
	LengthAsync(callback func(int))
	SetAsync(key Key, value Value, callback func(int))
}

// This is a thin wrapper over a thread-safe Map in order to provide additional
// funtionalities (such as async operations).
type concurrentMap struct {
	storage Map
}

func (cm *concurrentMap) UnderlyingStorage() map[Key]Value {
	return cm.storage.UnderlyingStorage()
}

func (cm *concurrentMap) Contains(key Key) bool {
	return cm.storage.Contains(key)
}

func (cm *concurrentMap) Clear() {
	cm.storage.Clear()
}

func (cm *concurrentMap) Delete(key Key) int {
	return cm.storage.Delete(key)
}

func (cm *concurrentMap) IsEmpty() bool {
	return cm.storage.IsEmpty()
}

func (cm *concurrentMap) Length() int {
	return cm.storage.Length()
}

func (cm *concurrentMap) Set(key Key, value Value) int {
	return cm.storage.Set(key, value)
}

func (cm *concurrentMap) ContainsAsync(key Key, callback func(bool)) {
	go func() {
		found := cm.Contains(key)
		callback(found)
	}()
}

func (cm *concurrentMap) UnderlyingStorageAsync(callback func(map[Key]Value)) {
	go func() {
		storage := cm.UnderlyingStorage()
		callback(storage)
	}()
}

func (cm *concurrentMap) ClearAsync(callback func()) {
	go func() {
		cm.Clear()
		callback()
	}()
}

func (cm *concurrentMap) DeleteAsync(key Key, callback func(int)) {
	go func() {
		result := cm.Delete(key)
		callback(result)
	}()
}

func (cm *concurrentMap) Get(key Key) (Value, bool) {
	return cm.storage.Get(key)
}

func (cm *concurrentMap) GetAsync(key Key, callback func(Value, bool)) {
	go func() {
		v, found := cm.Get(key)
		callback(v, found)
	}()
}

func (cm *concurrentMap) IsEmptyAsync(callback func(bool)) {
	go func() {
		isEmpty := cm.IsEmpty()
		callback(isEmpty)
	}()
}

func (cm *concurrentMap) LengthAsync(callback func(int)) {
	go func() {
		length := cm.Length()
		callback(length)
	}()
}

func (cm *concurrentMap) SetAsync(key Key, value Value, callback func(int)) {
	go func() {
		length := cm.Set(key, value)
		callback(length)
	}()
}

func newConcurrentMap(cm Map) ConcurrentMap {
	return &concurrentMap{storage: cm}
}
