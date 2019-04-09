package zipack

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	}
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

	// Test For file doesnt exist
	r2, err2 := mgr.GetReader(filepath.Join("default", "DOESNOTEXIST.txt"))
	if err2 == nil {
		t.Errorf("Exepcted error but call succeeded: %v", r2)
	}
	if err2 != os.ErrNotExist {
		t.Errorf("Expected os.ErrNotExist, Got %T", err2)
	}
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
	}
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

	// Test for file doesnt exist error
	r2, err2 := mgr.GetFileContents(filepath.Join("default", "DOESNOTEXIST.txt"))
	if err2 == nil {
		t.Errorf("Exepcted error but call succeeded: %v", r2)
	}
	if err2 != os.ErrNotExist {
		t.Errorf("Expected os.ErrNotExist, Got %T", err2)
	}
}

func BenchmarkGetReaderWCache(b *testing.B) {
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

func BenchmarkGetReaderWithDeleteCache(b *testing.B) {
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
