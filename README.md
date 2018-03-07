# gocontainer

[![Go Report Card](https://goreportcard.com/badge/github.com/protoman92/gocontainer)](https://goreportcard.com/report/github.com/protoman92/gocontainer)
[![Build Status](https://travis-ci.org/protoman92/gocontainer.svg?branch=master)](https://travis-ci.org/protoman92/gocontainer)
[![Coverage Status](https://coveralls.io/repos/github/protoman92/gocontainer/badge.svg?branch=master)](https://coveralls.io/github/protoman92/gocontainer?branch=master)

A collection for containers for go.

To use this package, run:

> go get github.com/protoman92/gocontainer

This package contains the following:

## gocollection: Collection implementations

Here we have **Collection** and **List** implementations that are somewhat like their Java equivalents. For thread-safety please use **ConcurrentList**.

## gomap: Key-value map implementation

Here we have **BasicMap** (light wrapper of a **map**), **ConcurrentMap** (thread-safe). There are 2 implementations of **ConcurrentMap**:

- **ChannelConcurrentMap**: Channel-based **ConcurrentMap** with each request type having its own channel and all coordination is done in a for loop within a goroutine.

- **LockConcurrentMap**: Simple mutex-dependent **ConcurrentMap**. This version should be faster than **ChannelConcurrentMap** based on benchmarks.