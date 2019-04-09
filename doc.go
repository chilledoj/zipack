// Package zipack provides a simple object to read files directly from a zip file and cache the results.
//
// The struct utilises the sync.Map instead of the built-in map as the cache will only increase and
// will not be deleting keys. This means that it should be safe for concurrent use.
//
// Example usage
//  mgr := NewManager("./testdata/sqlfiles.zip")
//
//  // Call Function
//  res, err := mgr.GetFileContents(filepath.Join("default", "selectDual.sql"))
//  if err != nil {
//    t.Errorf("Unexpected error: %v", err)
//  }
package zipack
