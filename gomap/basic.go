package gomap

import (
	"fmt"
)

// KeyStringer represents a key-string converter.
type KeyStringer func(Key) string

type basicMap struct {
	BasicMapParams
	storage map[string]Value
}

// BasicMapParams represents all the required parameters to build a BasicMap.
type BasicMapParams struct {
	InitialCap   uint
	ConvertKeyFn KeyStringer
}

func (b *basicMap) String() string {
	return fmt.Sprint(b.storage)
}

func (b *basicMap) convertKey(key Key) string {
	if b.ConvertKeyFn != nil {
		return b.ConvertKeyFn(key)
	}

	return fmt.Sprint(key)
}

// Returns a copy of the underlying storage.
func (b *basicMap) UnderlyingStorage() map[Key]Value {
	storage := make(map[Key]Value)

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

func (b *basicMap) Contains(key Key) bool {
	_, found := b.storage[b.convertKey(key)]
	return found
}

func (b *basicMap) Delete(key Key) bool {
	strKey := b.convertKey(key)
	prev := b.storage[strKey]
	delete(b.storage, strKey)
	return prev != nil
}

func (b *basicMap) Get(key Key) (Value, bool) {
	v, ok := b.storage[b.convertKey(key)]
	return v, ok
}

func (b *basicMap) IsEmpty() bool {
	return b.Length() == 0
}

func (b *basicMap) Length() int {
	return len(b.storage)
}

func (b *basicMap) Set(key Key, value Value) (Value, bool) {
	strKey := b.convertKey(key)
	prev := b.storage[strKey]
	b.storage[b.convertKey(key)] = value
	return prev, prev != nil
}

// NewBasicMap creates a new BasicMap.
func NewBasicMap(params BasicMapParams) Map {
	storage := make(map[string]Value)
	return &basicMap{BasicMapParams: params, storage: storage}
}

// NewDefaultBasicMap creates a new default BasicMap.
func NewDefaultBasicMap() Map {
	return NewBasicMap(BasicMapParams{})
}
