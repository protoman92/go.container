package gomap

import (
	"fmt"
)

// KeyStringer represents a key-string converter.
type KeyStringer func(interface{}) string

type basicMap struct {
	BasicMapParams
	storage map[string]interface{}
}

// BasicMapParams represents all the required parameters to build a BasicMap.
type BasicMapParams struct {
	InitialCap   uint
	ConvertKeyFn KeyStringer
}

func (b *basicMap) String() string {
	return fmt.Sprint(b.storage)
}

func (b *basicMap) convertKey(key interface{}) string {
	if b.ConvertKeyFn != nil {
		return b.ConvertKeyFn(key)
	}

	return fmt.Sprint(key)
}

// Returns a copy of the underlying storage.
func (b *basicMap) UnderlyingStorage() map[interface{}]interface{} {
	storage := make(map[interface{}]interface{})

	for key, value := range b.storage {
		storage[key] = value
	}

	return storage
}

func (b *basicMap) Clear() {
	for key := range b.storage {
		delete(b.storage, key)
	}
}

func (b *basicMap) Delete(key interface{}) int {
	delete(b.storage, b.convertKey(key))
	return len(b.storage)
}

func (b *basicMap) Get(key interface{}) (interface{}, bool) {
	v, ok := b.storage[b.convertKey(key)]
	return v, ok
}

func (b *basicMap) Length() int {
	return len(b.storage)
}

func (b *basicMap) Set(key interface{}, value interface{}) int {
	b.storage[b.convertKey(key)] = value
	return len(b.storage)
}

// NewBasicMap creates a new BasicMap.
func NewBasicMap(params BasicMapParams) Map {
	storage := make(map[string]interface{})
	return &basicMap{BasicMapParams: params, storage: storage}
}

// NewDefaultBasicMap creates a new default BasicMap.
func NewDefaultBasicMap() Map {
	return NewBasicMap(BasicMapParams{})
}
