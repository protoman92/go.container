package gomap

// Key represents the supported Map key type.
type Key interface{}

// Value represents the supported Map value type.
type Value interface{}

// Map represents a key-value storage. Thread-safety is not required.
type Map interface {
	UnderlyingStorage() map[Key]Value
	Clear()
	Contains(key Key) bool
	Delete(key Key) bool
	Get(key Key) (Value, bool)
	IsEmpty() bool
	Length() int

	// Set a key with a value, and return the previous value.
	Set(key Key, value Value) (Value, bool)
}
