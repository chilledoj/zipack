package zipack

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func BenchmarkGetReaderWCache(b *testing.B) {
	test := func(fileKey string) func(*testing.B) {
		return func(b *testing.B) {
			// SETUP
			mgr := NewManager(Options{ZipFileName: "./testdata/sqlfiles.zip"})
			for n := 0; n < b.N; n++ {
				// Call Function
				_, err := mgr.GetReader(fileKey)
				if err != nil {
					b.Errorf("Unexpected error: %v", err)
				}
			}
		}
	}
	b.Run("small file", test(filepath.Join("default", "selectDual.sql")))
	b.Run("larger file", test("largerFile.txt"))
}

func BenchmarkGetFileContentsWCache(b *testing.B) {
	test := func(fileKey string) func(*testing.B) {
		return func(b *testing.B) {
			// SETUP
			mgr := NewManager(Options{ZipFileName: "./testdata/sqlfiles.zip"})
			for n := 0; n < b.N; n++ {
				// Call Function
				_, err := mgr.GetFileContents(fileKey)
				if err != nil {
					b.Errorf("Unexpected error: %v", err)
				}
			}
		}
	}
	b.Run("small file", test(filepath.Join("default", "selectDual.sql")))
	b.Run("larger file", test("largerFile.txt"))
}
func BenchmarkOsOpenFileReadAllNonZipWSimilarCache(b *testing.B) {
	test := func(fileKey string) func(*testing.B) {
		return func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, err := getFileContentsNonZipWSimilarCache(fileKey)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	}
	b.Run("small file", test("./testdata/sqlfiles/default/selectDual.sql"))
	b.Run("larger file", test("./testdata/sqlfiles/largerFile.txt"))
}

func BenchmarkGetReaderWithDeleteCache(b *testing.B) {
	test := func(fileKey string) func(*testing.B) {
		return func(b *testing.B) {
			// SETUP
			mgr := NewManager(Options{ZipFileName: "./testdata/sqlfiles.zip"})
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				// Call Function
				_, err := mgr.GetReader(fileKey)
				if err != nil {
					b.Errorf("Unexpected error: %v", err)
				}
				b.StopTimer()
				mgr.cache.Delete(fileKey)
				b.StartTimer()
			}
		}
	}
	b.Run("small file", test(filepath.Join("default", "selectDual.sql")))
	b.Run("larger file", test("largerFile.txt"))
}

func BenchmarkGetFileContentsWithDeleteCache(b *testing.B) {
	test := func(fileKey string) func(*testing.B) {
		return func(b *testing.B) {
			// SETUP
			mgr := NewManager(Options{ZipFileName: "./testdata/sqlfiles.zip"})
			for n := 0; n < b.N; n++ {
				// Call Function
				_, err := mgr.GetFileContents(fileKey)
				if err != nil {
					b.Errorf("Unexpected error: %v", err)
				}
				b.StopTimer()
				mgr.cache.Delete(fileKey)
				b.StartTimer()
			}
		}
	}
	b.Run("small file", test(filepath.Join("default", "selectDual.sql")))
	b.Run("larger file", test("largerFile.txt"))

}
func BenchmarkOsOpenFileReadAllNonZipNoCache(b *testing.B) {
	test := func(fileKey string) func(*testing.B) {
		return func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				file, err := os.OpenFile(fileKey, os.O_RDONLY, 0666)
				if err != nil {
					b.Fatal(err)
				}
				_, err = ioutil.ReadAll(file)
				if err != nil {
					b.Fatal(err)
				}
				file.Close()
			}
		}
	}
	b.Run("small file", test("./testdata/sqlfiles/default/selectDual.sql"))
	b.Run("larger file", test("./testdata/sqlfiles/largerFile.txt"))
}

var cache sync.Map

func getFileContentsNonZipWSimilarCache(path string) ([]byte, error) {
	v, ok := cache.Load(path)
	if ok {
		return v.([]byte), nil
	}
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	cache.Store(path, buf)
	return buf, nil
}
