# gocollection

[![Go Report Card](https://goreportcard.com/badge/github.com/protoman92/gocollection)](https://goreportcard.com/report/github.com/protoman92/gocollection)
[![Build Status](https://travis-ci.org/protoman92/gocollection.svg?branch=master)](https://travis-ci.org/protoman92/gocollection)
[![Coverage Status](https://coveralls.io/repos/github/protoman92/gocollection/badge.svg?branch=master)](https://coveralls.io/github/protoman92/gocollection?branch=master)

A collection for collections for go.

To use this package, run:

> go get github.com/protoman92/gocollection

This package contains the following:

## gomap: Key-value map implementation

Here we have **BasicMap** (light wrapper of a **map**), **ConcurrentMap** (thread-safe). There are 2 implementations of **ConcurrentMap**:

- **ChannelConcurrentMap**: Channel-based **ConcurrentMap** with each request type having its own channel and all coordination is done in a for loop within a goroutine.

- **LockConcurrentMap**: Simple mutex-dependent **ConcurrentMap**. This version should be faster than **ChannelConcurrentMap** based on benchmarks.