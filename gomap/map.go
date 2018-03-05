package gomap

// Map represents a key-value storage. Thread-safety is not required.
type Map interface {
	UnderlyingStorage() map[interface{}]interface{}
	Clear()
	Get(key interface{}) (interface{}, bool)
	Length() int

	// Delete a key and return the new length.
	Delete(key interface{}) int

	// Set a key with a value, and return the new length.
	Set(key interface{}, value interface{}) int
}
