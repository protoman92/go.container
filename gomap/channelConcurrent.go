package gomap

type containsRequest struct {
	key     Key
	foundCh chan<- bool
}

type deleteRequest struct {
	key   Key
	lenCh chan<- bool
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

type channelConcurrentMap struct {
	storage         Map
	accessMapCh     chan chan Map
	accessStorageCh chan chan map[Key]Value
	clearCh         chan chan interface{}
	containsCh      chan *containsRequest
	deleteCh        chan *deleteRequest
	lenCh           chan chan int
	getCh           chan *getRequest
	setCh           chan *setRequest
}

// This operation blocks until the underlying storage is received.
func (ccm *channelConcurrentMap) UnderlyingStorage() map[Key]Value {
	accessCh := make(chan map[Key]Value, 0)
	ccm.accessStorageCh <- accessCh
	return <-accessCh
}

// This operation blocks until some result is received.
func (ccm *channelConcurrentMap) Clear() {
	requestCh := make(chan interface{}, 0)
	ccm.clearCh <- requestCh
	<-requestCh
}

// This operation blocks until a value is received.
func (ccm *channelConcurrentMap) Contains(key Key) bool {
	foundCh := make(chan bool, 0)
	ccm.containsCh <- &containsRequest{key: key, foundCh: foundCh}
	return <-foundCh
}

// This operation blocks until some value is received.
func (ccm *channelConcurrentMap) Delete(key Key) bool {
	lenCh := make(chan bool, 0)
	ccm.deleteCh <- &deleteRequest{key: key, lenCh: lenCh}
	return <-lenCh
}

// This operaton blocks until some value is received.
func (ccm *channelConcurrentMap) Get(key Key) (Value, bool) {
	valueCh := make(chan *getResult, 0)
	ccm.getCh <- &getRequest{key: key, valueCh: valueCh}
	result := <-valueCh
	return result.element, result.found
}

// This operaton blocks until some value is received.
func (ccm *channelConcurrentMap) IsEmpty() bool {
	return ccm.Length() == 0
}

// This operaton blocks until some value is received.
func (ccm *channelConcurrentMap) Length() int {
	requestCh := make(chan int, 0)
	ccm.lenCh <- requestCh
	return <-requestCh
}

// This operaton blocks until some value is received.
func (ccm *channelConcurrentMap) Set(key Key, value Value) (Value, bool) {
	lenCh := make(chan *setResult, 0)
	ccm.setCh <- &setRequest{key: key, value: value, lenCh: lenCh}
	result := <-lenCh
	return result.element, result.found
}

// This operation blocks until the underlying Map is received.
func (ccm *channelConcurrentMap) UnderlyingMap() Map {
	accessCh := make(chan Map, 0)
	ccm.accessMapCh <- accessCh
	return <-accessCh
}

func (ccm *channelConcurrentMap) loopMap() {
	for {
		select {
		case ar := <-ccm.accessMapCh:
			ar <- ccm.storage

		case ar := <-ccm.accessStorageCh:
			ar <- ccm.storage.UnderlyingStorage()

		case cr := <-ccm.clearCh:
			ccm.storage.Clear()
			cr <- true

		case cr := <-ccm.containsCh:
			cr.foundCh <- ccm.storage.Contains(cr.key)

		case dr := <-ccm.deleteCh:
			dr.lenCh <- ccm.storage.Delete(dr.key)

		case lr := <-ccm.lenCh:
			lr <- ccm.storage.Length()

		case gr := <-ccm.getCh:
			element, found := ccm.storage.Get(gr.key)
			gr.valueCh <- &getResult{element: element, found: found}

		case sr := <-ccm.setCh:
			element, found := ccm.storage.Set(sr.key, sr.value)
			sr.lenCh <- &setResult{element: element, found: found}
		}
	}
}

// NewChannelConcurrentMap returns a new channel-based ConcurrentMap.
func NewChannelConcurrentMap(storage Map) ConcurrentMap {
	cm := &channelConcurrentMap{
		storage:         storage,
		accessMapCh:     make(chan chan Map, 0),
		accessStorageCh: make(chan chan map[Key]Value, 0),
		clearCh:         make(chan chan interface{}, 0),
		containsCh:      make(chan *containsRequest, 0),
		deleteCh:        make(chan *deleteRequest, 0),
		lenCh:           make(chan chan int, 0),
		getCh:           make(chan *getRequest, 0),
		setCh:           make(chan *setRequest, 0),
	}

	go cm.loopMap()
	return newConcurrentMap(cm)
}
