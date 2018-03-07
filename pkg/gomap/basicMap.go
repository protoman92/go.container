package gomap

import (
	"fmt"
)

type basicMap struct {
	BasicMapParams
	storage map[interface{}]interface{}
}

// BasicMapParams represents all the required parameters to build a BasicMap.
type BasicMapParams struct {
	InitialCap uint
}

func (b *basicMap) String() string {
	return fmt.Sprint(b.storage)
}

func (b *basicMap) Clear() {
	for key := range b.storage {
		delete(b.storage, key)
	}
}

func (b *basicMap) Contains(key interface{}) bool {
	_, found := b.storage[key]
	return found
}

func (b *basicMap) Delete(key interface{}) (interface{}, bool) {
	prev := b.storage[key]
	delete(b.storage, key)
	return prev, prev != nil
}

func (b *basicMap) Get(key interface{}) (interface{}, bool) {
	v, ok := b.storage[key]
	return v, ok
}

func (b *basicMap) Length() int {
	return len(b.storage)
}

func (b *basicMap) Keys() []interface{} {
	keys := make([]interface{}, 0)

	for key := range b.storage {
		keys = append(keys, key)
	}

	return keys
}

func (b *basicMap) Set(key interface{}, value interface{}) (interface{}, bool) {
	prev := b.storage[key]
	b.storage[key] = value
	return prev, prev != nil
}

// NewBasicMap creates a new BasicMap.
func NewBasicMap(params BasicMapParams) Map {
	storage := make(map[interface{}]interface{})
	return &basicMap{BasicMapParams: params, storage: storage}
}

// NewDefaultBasicMap creates a new default BasicMap.
func NewDefaultBasicMap() Map {
	return NewBasicMap(BasicMapParams{})
}
