package gomap

import (
	"sync"
	"testing"
	"time"
)

func TestUnsupportedRequestShouldPanic(t *testing.T) {
	/// Setup
	bm := NewDefaultBasicMap()
	requestCh := make(chan interface{}, 1)
	ccm := &channelConcurrentMap{storage: bm, requestCh: requestCh}
	var e string
	var mutex sync.RWMutex

	accessError := func() string {
		mutex.RLock()
		defer mutex.RUnlock()
		return e
	}

	setError := func(err interface{}) {
		mutex.Lock()
		defer mutex.Unlock()
		msg := err.(string)
		e = msg
	}

	go func() {
		defer func() {
			if e := recover(); e == nil {
				t.Errorf("Should have panicked")
			} else {
				setError(e)
			}
		}()

		ccm.loopMap()
	}()

	/// When
	requestCh <- true

	/// Then
	time.Sleep(time.Millisecond)

	if err := accessError(); len(err) == 0 {
		t.Errorf("Should have panicked")
	}
}

func TestChannelConncurrentMapCloseRequest(t *testing.T) {
	/// Setup
	doneCh := make(chan interface{}, 1)
	timeout := time.After(time.Second)
	bm := NewBasicMap(BasicMapParams{})
	cm := NewChannelConcurrentMap(bm)
	key := "Key"

	/// When
	cm.Close()

	go func() {
		defer func() {
			if e := recover(); e != nil {
				doneCh <- true
			}
		}()

		cm.Set(key, "Value")
	}()

	time.Sleep(time.Millisecond)

	/// Then
	select {
	case <-doneCh:
		break

	case <-timeout:
		t.Errorf("Should have quit loop")
	}
}
