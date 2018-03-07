// Package gomap provides thread-safe ConcurrentMap implementations, of which
// there are 2: lock-based and channel-based ConcurrentMap. To construct a
// ConcurrentMap, we need to supply a non-thread safe Map implementation, such
// as BasicMap. Note that the lock variant is faster than the channel version.
package gomap
