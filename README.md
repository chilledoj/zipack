# zipack

[![GoDoc](https://godoc.org/github.com/chilledoj/zipack?status.svg)](http://godoc.org/github.com/chilledoj/zipack)
[![Go Report Card](https://goreportcard.com/badge/github.com/chilledoj/zipack)](https://goreportcard.com/report/github.com/chilledoj/zipack)
[![goreadme](https://goreadme.herokuapp.com/badge/chilledoj/zipack.svg)](https://goreadme.herokuapp.com)

Package zipack provides a simple object to read files directly from a zip file and cache the results.
The code assumes that a folder has been zip'd and thus the files reside within the folder

The struct utilises the sync.Map instead of the built-in map as the cache will only increase and
will not be deleting keys. This means that it should be safe for concurrent use.


---

Created by [goreadme](https://github.com/apps/goreadme)
