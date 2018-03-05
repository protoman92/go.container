package gomap

import "sync"

type containsRequest struct {
	key     Key
	foundCh chan<- bool
}

type deleteRequest struct {
	key   Key
	lenCh chan<- int
}

type getResult struct {
	element Value
	found   bool
}

type getRequest struct {
	key     Key
	valueCh chan<- *getResult
}

type setRequest struct {
	key   Key
	value Value
	lenCh chan<- int
}

// ConcurrentMap represents a thread-safe Map.
type ConcurrentMap interface {
	Map
	UnderlyingMap() Map
	ClearAsync(callback func())
	ContainsAsync(key Key, callback func(bool))
	DeleteAsync(key Key, callback func(int))
	GetAsync(key Key, callback func(Value, bool))
	IsEmptyAsync(callback func(bool))
	LengthAsync(callback func(int))
	SetAsync(key Key, value Value, callback func(int))
}

// This is a wrapper over a Map that provides thread-safe operations.
type concurrentMap struct {
	mutex      sync.RWMutex
	storage    Map
	clearCh    chan chan interface{}
	containsCh chan *containsRequest
	deleteCh   chan *deleteRequest
	lenCh      chan chan int
	getCh      chan *getRequest
	setCh      chan *setRequest
}

func (cm *concurrentMap) UnderlyingMap() Map {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.storage
}

func (cm *concurrentMap) UnderlyingStorage() map[Key]Value {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.storage.UnderlyingStorage()
}

func (cm *concurrentMap) Clear() {
	requestCh := make(chan interface{})
	cm.clearCh <- requestCh
	<-requestCh
}

func (cm *concurrentMap) ClearAsync(callback func()) {
	go func() {
		cm.Clear()
		callback()
	}()
}

func (cm *concurrentMap) Contains(key Key) bool {
	foundCh := make(chan bool)
	cm.containsCh <- &containsRequest{key: key, foundCh: foundCh}
	return <-foundCh
}

func (cm *concurrentMap) ContainsAsync(key Key, callback func(bool)) {
	go func() {
		found := cm.Contains(key)
		callback(found)
	}()
}

func (cm *concurrentMap) Delete(key Key) int {
	lenCh := make(chan int)
	cm.deleteCh <- &deleteRequest{key: key, lenCh: lenCh}
	return <-lenCh
}

func (cm *concurrentMap) DeleteAsync(key Key, callback func(int)) {
	go func() {
		result := cm.Delete(key)
		callback(result)
	}()
}

func (cm *concurrentMap) Get(key Key) (Value, bool) {
	valueCh := make(chan *getResult, 0)
	cm.getCh <- &getRequest{key: key, valueCh: valueCh}
	result := <-valueCh
	return result.element, result.found
}

func (cm *concurrentMap) GetAsync(key Key, callback func(Value, bool)) {
	go func() {
		v, found := cm.Get(key)
		callback(v, found)
	}()
}

func (cm *concurrentMap) IsEmpty() bool {
	return cm.Length() == 0
}

func (cm *concurrentMap) IsEmptyAsync(callback func(bool)) {
	go func() {
		isEmpty := cm.IsEmpty()
		callback(isEmpty)
	}()
}

func (cm *concurrentMap) Length() int {
	requestCh := make(chan int)
	cm.lenCh <- requestCh
	return <-requestCh
}

func (cm *concurrentMap) LengthAsync(callback func(int)) {
	go func() {
		length := cm.Length()
		callback(length)
	}()
}

func (cm *concurrentMap) Set(key Key, value Value) int {
	lenCh := make(chan int, 0)
	cm.setCh <- &setRequest{key: key, value: value, lenCh: lenCh}
	return <-lenCh
}

func (cm *concurrentMap) SetAsync(key Key, value Value, callback func(int)) {
	go func() {
		length := cm.Set(key, value)
		callback(length)
	}()
}

func (cm *concurrentMap) loopMap() {
	for {
		select {
		case cr := <-cm.clearCh:
			cm.storage.Clear()
			cr <- true

		case cr := <-cm.containsCh:
			cr.foundCh <- cm.storage.Contains(cr.key)

		case dr := <-cm.deleteCh:
			dr.lenCh <- cm.storage.Delete(dr.key)

		case lr := <-cm.lenCh:
			lr <- cm.storage.Length()

		case gr := <-cm.getCh:
			element, found := cm.storage.Get(gr.key)
			gr.valueCh <- &getResult{element: element, found: found}

		case sr := <-cm.setCh:
			sr.lenCh <- cm.storage.Set(sr.key, sr.value)
		}
	}
}

// NewConcurrentMap returns a new ConcurrentMap.
func NewConcurrentMap(storage Map) ConcurrentMap {
	cm := &concurrentMap{
		storage:    storage,
		clearCh:    make(chan chan interface{}, 0),
		containsCh: make(chan *containsRequest, 0),
		deleteCh:   make(chan *deleteRequest, 0),
		lenCh:      make(chan chan int, 0),
		getCh:      make(chan *getRequest, 0),
		setCh:      make(chan *setRequest, 0),
	}

	go cm.loopMap()
	return cm
}

// NewDefaultBasicConcurrentMap returns a new ConcurrentMap backed by a default
// BasicMap.
func NewDefaultBasicConcurrentMap() ConcurrentMap {
	storage := NewDefaultBasicMap()
	return NewConcurrentMap(storage)
}
