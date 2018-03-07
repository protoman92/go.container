package gocollection

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type ConcurrentOpsParams struct {
	concurrentCollection Collection
	log                  bool
	keys                 []interface{}
	opSleepDuration      func() time.Duration
}

func setupConcurrentCollectionOps(params *ConcurrentOpsParams) func() *sync.WaitGroup {
	/// Setup
	cl := params.concurrentCollection
	keys := params.keys

	for ix := range keys {
		keys[ix] = strconv.Itoa(ix)
	}

	waitGroup := &sync.WaitGroup{}
	mutex := sync.RWMutex{}

	accessWaitGroup := func() *sync.WaitGroup {
		mutex.Lock()
		defer mutex.Unlock()
		return waitGroup
	}

	/// When & Then
	// Modify
	go func() {
		for _, key := range keys {
			accessWaitGroup().Add(1)

			go func(key interface{}) {
				time.Sleep(params.opSleepDuration())

				if added := cl.Add(key); params.log {
					fmt.Printf("Added %v, added count: %v\n", key, added)
				}

				accessWaitGroup().Done()
			}(key)
		}
	}()

	go func() {
		for i := 0; i < len(keys); i++ {
			accessWaitGroup().Add(1)

			go func() {
				time.Sleep(params.opSleepDuration())

				if added := cl.AddAll(keys...); params.log {
					fmt.Printf("Added %d elements\n", added)
				}

				accessWaitGroup().Done()
			}()
		}
	}()

	go func() {
		for _, key := range keys {
			accessWaitGroup().Add(1)

			go func(key interface{}) {
				time.Sleep(params.opSleepDuration())

				if removed := cl.Remove(key); params.log {
					fmt.Printf("Deleted %v, removed count: %d\n", key, removed)
				}

				accessWaitGroup().Done()
			}(key)
		}
	}()

	go func() {
		for i := 0; i < len(keys); i++ {
			accessWaitGroup().Add(1)

			go func() {
				time.Sleep(params.opSleepDuration())
				cl.Clear()

				if params.log {
					fmt.Printf("Cleared all contents\n")
				}

				accessWaitGroup().Done()
			}()
		}
	}()

	// Get
	go func() {
		for _, key := range keys {
			accessWaitGroup().Add(1)

			go func(key interface{}) {
				time.Sleep(params.opSleepDuration())

				if found := cl.Contains(key); params.log {
					fmt.Printf("Contains element %v: %t\n", key, found)
				}

				accessWaitGroup().Done()
			}(key)
		}
	}()

	go func() {
		for i := 0; i < len(keys); i++ {
			accessWaitGroup().Add(1)

			go func() {
				time.Sleep(params.opSleepDuration())

				if len := cl.Length(); params.log {
					fmt.Printf("Current length: %d\n", len)
				}

				accessWaitGroup().Done()
			}()
		}
	}()

	return accessWaitGroup
}
