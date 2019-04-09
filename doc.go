// Package zipack provides a simple object to read files directly from a zip file and cache the results.
//
// The struct utilises the sync.Map instead of the built-in map as the cache will only increase and
// will not be deleting keys. This means that it should be safe for concurrent use.
//
package zipack
