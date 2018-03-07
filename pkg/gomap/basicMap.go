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

func (b *basicMap) Clear() {
	for key := range b.storage {
		delete(b.storage, key)
	}
}

func (b *basicMap) Contains(key interface{}) bool {
	_, found := b.storage[b.convertKey(key)]
	return found
}

func (b *basicMap) Delete(key interface{}) bool {
	strKey := b.convertKey(key)
	prev := b.storage[strKey]
	delete(b.storage, strKey)
	return prev != nil
}

func (b *basicMap) Get(key interface{}) (interface{}, bool) {
	v, ok := b.storage[b.convertKey(key)]
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
	strKey := b.convertKey(key)
	prev := b.storage[strKey]
	b.storage[strKey] = value
	return prev, prev != nil
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
