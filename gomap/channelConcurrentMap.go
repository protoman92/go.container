package gomap

import (
	"fmt"
	"reflect"
)

type clearRequest struct {
	doneCh chan<- interface{}
}

type containsRequest struct {
	key     Key
	foundCh chan<- bool
}

type deleteRequest struct {
	key   Key
	lenCh chan<- bool
}

type lenRequest struct {
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

type setResult struct {
	element Value
	found   bool
}

type setRequest struct {
	key   Key
	value Value
	lenCh chan<- *setResult
}

type stringRequest struct {
	strCh chan<- string
}

type channelConcurrentMap struct {
	storage   Map
	requestCh chan interface{}
}

func (ccm *channelConcurrentMap) String() string {
	strCh := make(chan string, 0)
	ccm.requestCh <- &stringRequest{strCh: strCh}
	return <-strCh
}

// This operation blocks until some result is received.
func (ccm *channelConcurrentMap) Clear() {
	requestCh := make(chan interface{}, 0)
	ccm.requestCh <- &clearRequest{doneCh: requestCh}
	<-requestCh
}

// This operation blocks until a value is received.
func (ccm *channelConcurrentMap) Contains(key Key) bool {
	foundCh := make(chan bool, 0)
	ccm.requestCh <- &containsRequest{key: key, foundCh: foundCh}
	return <-foundCh
}

// This operation blocks until some value is received.
func (ccm *channelConcurrentMap) Delete(key Key) bool {
	lenCh := make(chan bool, 0)
	ccm.requestCh <- &deleteRequest{key: key, lenCh: lenCh}
	return <-lenCh
}

// This operaton blocks until some value is received.
func (ccm *channelConcurrentMap) Get(key Key) (Value, bool) {
	valueCh := make(chan *getResult, 0)
	ccm.requestCh <- &getRequest{key: key, valueCh: valueCh}
	result := <-valueCh
	return result.element, result.found
}

// This operaton blocks until some value is received.
func (ccm *channelConcurrentMap) Length() int {
	requestCh := make(chan int, 0)
	ccm.requestCh <- &lenRequest{lenCh: requestCh}
	return <-requestCh
}

// This operaton blocks until some value is received.
func (ccm *channelConcurrentMap) Set(key Key, value Value) (Value, bool) {
	lenCh := make(chan *setResult, 0)
	ccm.requestCh <- &setRequest{key: key, value: value, lenCh: lenCh}
	result := <-lenCh
	return result.element, result.found
}

func (ccm *channelConcurrentMap) loopMap() {
	for {
		select {
		case request := <-ccm.requestCh:
			switch request := request.(type) {
			case *clearRequest:
				ccm.storage.Clear()
				request.doneCh <- true

			case *containsRequest:
				request.foundCh <- ccm.storage.Contains(request.key)

			case *deleteRequest:
				request.lenCh <- ccm.storage.Delete(request.key)

			case *lenRequest:
				request.lenCh <- ccm.storage.Length()

			case *getRequest:
				element, found := ccm.storage.Get(request.key)
				request.valueCh <- &getResult{element: element, found: found}

			case *setRequest:
				element, found := ccm.storage.Set(request.key, request.value)
				request.lenCh <- &setResult{element: element, found: found}

			case *stringRequest:
				request.strCh <- fmt.Sprint(ccm.storage)

			default:
				panic(fmt.Sprintf("Unrecognized req type %v", reflect.TypeOf(request)))
			}
		}
	}
}

// NewChannelConcurrentMap returns a new channel-based ConcurrentMap.
func NewChannelConcurrentMap(storage Map) ConcurrentMap {
	cm := &channelConcurrentMap{
		storage:   storage,
		requestCh: make(chan interface{}, 1),
	}

	go cm.loopMap()
	return newConcurrentMap(cm)
}
