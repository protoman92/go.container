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
	Get(key Key) (Value, bool)
	IsEmpty() bool
	Length() int

	// Delete a key and return the new length.
	Delete(key Key) int

	// Set a key with a value, and return the new length.
	Set(key Key, value Value) int
}
