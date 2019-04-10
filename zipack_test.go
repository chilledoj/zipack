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
	mgr := NewManager(Options{ZipFileName: "./testdata/sqlfiles.zip"})

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

	t.Run("File does not exist", func(t2 *testing.T) {
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
		m := NewManager(Options{ZipFileName: "./testdata/sqlfiles.zip"})
		_, err := m.GetReader("doesn't matter")
		if err == nil {
			t2.Errorf("Get worked without error")
		}
	})
}

func TestGetFileContents(t *testing.T) {
	// SETUP
	mgr := NewManager(Options{ZipFileName: "./testdata/sqlfiles.zip"})

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

func TestPreloadPaths(t *testing.T) {
	mgr := NewManager(Options{
		ZipFileName: "./testdata/sqlfiles.zip",
		PreloadPaths: []string{
			filepath.Join("simple", "aa.txt"),
			filepath.Join("simple", "bb.txt"),
			filepath.Join("simple", "cc.txt"),
		},
	})

	tests := []struct {
		Path    string
		Present bool
	}{
		{Path: filepath.Join("simple", "aa.txt"), Present: true},
		{Path: filepath.Join("simple", "bb.txt"), Present: true},
		{Path: filepath.Join("simple", "cc.txt"), Present: true},
		{Path: filepath.Join("simple", "dd.txt"), Present: false},
		{Path: filepath.Join("simple", "ee.txt"), Present: false},
		{Path: filepath.Join("simple", "ff.txt"), Present: false},
	}
	for _, test := range tests {
		_, ok := mgr.cache.Load(test.Path)
		if ok != test.Present {
			t.Errorf("Path %s preload: Act(%v) Exp(%v)", test.Path, ok, test.Present)
		}
	}
}

func TestReadAndStoreFromZip(t *testing.T) {
	t.Run("zip file exists", func(t *testing.T) {
		mgr := NewManager(Options{ZipFileName: "./testdata/sqlfiles.zip"})
		_, err := mgr.readAndStoreFromZip("largerFile.txt")
		if err != nil {
			t.Errorf("Did not expect error: %v", err)
		}
	})
	t.Run("zip file not exists", func(t *testing.T) {
		mgr := NewManager(Options{ZipFileName: "./testdata/fail.zip"})
		_, err := mgr.readAndStoreFromZip("largerFile.txt")
		if err == nil {
			t.Error("Expected error")
		}
	})
}

func ExampleManager_GetFileContents() {
	mgr := NewManager(Options{ZipFileName: "./testdata/sqlfiles.zip"})
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
