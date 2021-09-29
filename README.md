# zipack

[![GoDoc](https://godoc.org/github.com/chilledoj/zipack?status.svg)](http://godoc.org/github.com/chilledoj/zipack)
[![Go Report Card](https://goreportcard.com/badge/github.com/chilledoj/zipack)](https://goreportcard.com/report/github.com/chilledoj/zipack)

```shell
go get -u github.com/chilledoj/zipack
```

Package zipack provides a simple object to read files directly from a zip file and cache the results.

The struct utilises a [sync.Map](https://pkg.go.dev/sync#Map) instead of the built-in map as the cache will only increase and will not be 
deleting keys. 
This also means that it should be safe for concurrent use.

Care should be taken when creating a zip file. For instance, the `sqlfiles.zip` file in the `testdata` folder was 
created on a Mac by selecting the folder in the finder, and then selecting the compress option on the context menu. 
This creates the zip file with the folder included in the file structure. Therefore all of your paths must
include this folder. See below and the tests for examples.

## Example
```go
mgr := NewManager(Options{ZipFileName: "./testdata/sqlfiles.zip"})
fileKey := filepath.Join("sqlfiles","default", "selectDual.sql")
res, err := mgr.GetFileContents(fileKey)
if err != nil {
    panic(fmt.Errorf("Unexpected error: %v", err))
}
fmt.Printf("%s", res)
// Output:
// select
//   *
// from
//   DUAL
```