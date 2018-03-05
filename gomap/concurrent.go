package gomap

import "sync"

type deleteRequest struct {
	key   interface{}
	lenCh chan<- int
}

type getResult struct {
	element interface{}
	found   bool
}

type getRequest struct {
	key     interface{}
	valueCh chan<- *getResult
}

type setRequest struct {
	key   interface{}
	value interface{}
	lenCh chan<- int
}

// ConcurrentMap represents a thread-safe Map.
type ConcurrentMap interface {
	Map
	UnderlyingMap() Map
	ClearAsync(callback func())
	DeleteAsync(key interface{}, callback func(int))
	GetAsync(key interface{}, callback func(interface{}, bool))
	LengthAsync(callback func(int))
	SetAsync(key interface{}, value interface{}, callback func(int))
}

// This is a wrapper over a Map that provides thread-safe operations.
type concurrentMap struct {
	mutex    sync.RWMutex
	storage  Map
	clearCh  chan chan interface{}
	deleteCh chan *deleteRequest
	lenCh    chan chan int
	getCh    chan *getRequest
	setCh    chan *setRequest
}

func (cm *concurrentMap) UnderlyingMap() Map {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.storage
}

func (cm *concurrentMap) UnderlyingStorage() map[interface{}]interface{} {
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

func (cm *concurrentMap) Delete(key interface{}) int {
	lenCh := make(chan int)
	cm.deleteCh <- &deleteRequest{key: key, lenCh: lenCh}
	return <-lenCh
}

func (cm *concurrentMap) DeleteAsync(key interface{}, callback func(int)) {
	go func() {
		result := cm.Delete(key)
		callback(result)
	}()
}

func (cm *concurrentMap) Get(key interface{}) (interface{}, bool) {
	valueCh := make(chan *getResult, 0)
	cm.getCh <- &getRequest{key: key, valueCh: valueCh}
	result := <-valueCh
	return result.element, result.found
}

func (cm *concurrentMap) GetAsync(key interface{}, callback func(interface{}, bool)) {
	go func() {
		v, found := cm.Get(key)
		callback(v, found)
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

func (cm *concurrentMap) Set(key interface{}, value interface{}) int {
	lenCh := make(chan int, 0)
	cm.setCh <- &setRequest{key: key, value: value, lenCh: lenCh}
	return <-lenCh
}

func (cm *concurrentMap) SetAsync(key interface{}, value interface{}, callback func(int)) {
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
		storage:  storage,
		clearCh:  make(chan chan interface{}, 0),
		deleteCh: make(chan *deleteRequest, 0),
		lenCh:    make(chan chan int, 0),
		getCh:    make(chan *getRequest, 0),
		setCh:    make(chan *setRequest, 0),
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
