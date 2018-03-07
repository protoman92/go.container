package gomap

// Map represents a key-value storage. Thread-safety is not required.
type Map interface {
	Clear()
	Contains(key interface{}) bool
	Delete(key interface{}) (interface{}, bool)
	Get(key interface{}) (interface{}, bool)
	Keys() []interface{}
	Length() int

	// Set a key with a value, and return the previous value.
	Set(key interface{}, value interface{}) (interface{}, bool)
}
