package gocollection

import (
	"sync"
)

type lockConcurrentCollection struct {
	mutex   *sync.RWMutex
	storage Collection
}

func (c *lockConcurrentCollection) Add(element interface{}) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.storage.Add(element)
}

func (c *lockConcurrentCollection) AddAll(elements ...interface{}) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.storage.AddAll(elements...)
}

func (c *lockConcurrentCollection) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.storage.Clear()
}

func (c *lockConcurrentCollection) Contains(element interface{}) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.storage.Contains(element)
}

func (c *lockConcurrentCollection) ContainsAll(elements ...interface{}) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.storage.ContainsAll(elements...)
}

func (c *lockConcurrentCollection) GetAllFunc(selector func(interface{}) bool) []interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.storage.GetAllFunc(selector)
}

func (c *lockConcurrentCollection) Length() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.storage.Length()
}

func (c *lockConcurrentCollection) Remove(element interface{}) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.storage.Remove(element)
}

func (c *lockConcurrentCollection) RemoveAll(elements ...interface{}) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.storage.RemoveAll(elements...)
}
