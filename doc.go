// Package zipack provides a simple object to read files directly from a zip zipfile and cache the results.
//
// The struct utilises the sync.Map instead of the built-in map as the cache will only increase and
// will not be deleting keys. This means that it should be safe for concurrent use.
//
// Care should be taken when creating a zip zipfile. For instance, the `sqlfiles.zip` zipfile in the `testdata` folder
// was created on a Mac by selecting the folder in the finder, and then selecting the compress option on the context
// menu. This creates the zip zipfile with the folder included in the zipfile structure. Therefore all of your paths must
// include this folder. See the tests for an example.
package zipack
