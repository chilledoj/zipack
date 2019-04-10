package zipack

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestGetReader(t *testing.T) {
	// SETUP
	mgr := NewManager("./testdata/sqlfiles.zip")

	// Call Function
	r, err := mgr.GetReader(filepath.Join("default", "selectDual.sql"))
	if err != nil {
		t.Errorf("Got unexepcted error: %v", err)
	}

	res, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Got unexepcted error: %v", err)
	}

	file, err := os.OpenFile("./testdata/sqlfiles/default/selectDual.sql", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	filedata, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != len(filedata) {
		t.Errorf("Length not the same : Act (%d) != Exp (%d)", len(res), len(filedata))
	} else {
		same := true
		for i := 0; i < len(res); i++ {
			if res[i] != filedata[i] {
				same = false
				break
			}
		}
		if !same {
			t.Errorf("Act (%s) != Exp (%s)", res, filedata)
		}
	}

	t.Run("File does not extist", func(t2 *testing.T) {
		r2, err2 := mgr.GetReader(filepath.Join("default", "DOESNOTEXIST.txt"))
		if err2 == nil {
			t2.Errorf("Exepcted error but call succeeded: %v", r2)
		}
		if err2 != os.ErrNotExist {
			t2.Errorf("Expected os.ErrNotExist, Got %T", err2)
		}
	})

	t.Run("Use cache", func(t2 *testing.T) {
		_, err3 := mgr.GetReader(filepath.Join("default", "selectDual.sql"))
		if err3 != nil {
			t2.Errorf("Unexepcted error: %v", err3)
		}
	})

	t.Run("Nonzip file", func(t2 *testing.T) {
		m := NewManager("./testdata/sqlfiles/default/selectDual.sql")
		_, err := m.GetReader("doesn't matter")
		if err == nil {
			t2.Errorf("Get worked without error")
		}
	})
}

func TestGetFileContents(t *testing.T) {
	// SETUP
	mgr := NewManager("./testdata/sqlfiles.zip")

	// Call Function
	res, err := mgr.GetFileContents(filepath.Join("default", "selectDual.sql"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check Against Fixtures
	file, err := os.OpenFile("./testdata/sqlfiles/default/selectDual.sql", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	filedata, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != len(filedata) {
		t.Errorf("Length not the same : Act (%d) != Exp (%d)", len(res), len(filedata))
	} else {
		same := true
		for i := 0; i < len(res); i++ {
			if res[i] != filedata[i] {
				same = false
				break
			}
		}
		if !same {
			t.Errorf("Act (%s) != Exp (%s)", res, filedata)
		}
	}

	t.Run("File does not exist", func(t2 *testing.T) {
		r2, err2 := mgr.GetFileContents(filepath.Join("default", "DOESNOTEXIST.txt"))
		if err2 == nil {
			t2.Errorf("Exepcted error but call succeeded: %v", r2)
		}
		if err2 != os.ErrNotExist {
			t2.Errorf("Expected os.ErrNotExist, Got %T", err2)
		}
	})

	t.Run("Use cache", func(t2 *testing.T) {
		_, err3 := mgr.GetFileContents(filepath.Join("default", "selectDual.sql"))
		if err3 != nil {
			t2.Errorf("Unexepcted error: %v", err3)
		}
	})

}

func BenchmarkGetReaderWCache(b *testing.B) {
	// SETUP
	mgr := NewManager("./testdata/sqlfiles.zip")
	fileKey := filepath.Join("default", "selectDual.sql")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// Call Function
		_, err := mgr.GetReader(fileKey)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkGetFileContentsWCache(b *testing.B) {
	// SETUP
	mgr := NewManager("./testdata/sqlfiles.zip")
	fileKey := filepath.Join("default", "selectDual.sql")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// Call Function
		_, err := mgr.GetFileContents(fileKey)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}
func BenchmarkOsOpenFileReadAllNonZipWSimilarCache(b *testing.B) {
	fileKey := "./testdata/sqlfiles/default/selectDual.sql"
	for n := 0; n < b.N; n++ {
		_, err := getFileContentsNonZipWSimilarCache(fileKey)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetReaderWithDeleteCache(b *testing.B) {
	// SETUP
	mgr := NewManager("./testdata/sqlfiles.zip")
	fileKey := filepath.Join("default", "selectDual.sql")
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

func BenchmarkGetFileContentsWithDeleteCache(b *testing.B) {
	// SETUP
	mgr := NewManager("./testdata/sqlfiles.zip")
	fileKey := filepath.Join("default", "selectDual.sql")
	b.ResetTimer()
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
func BenchmarkOsOpenFileReadAllNonZipNoCache(b *testing.B) {
	fileKey := "./testdata/sqlfiles/default/selectDual.sql"
	b.ResetTimer()
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

func ExampleManager_GetFileContents() {
	mgr := NewManager("./testdata/sqlfiles.zip")
	fileKey := filepath.Join("default", "selectDual.sql")
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
}
