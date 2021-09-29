package zipack

import (
	"bytes"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"
)

type tstFileInfo struct {
	Msg []byte
	Dir bool
}

func (t tstFileInfo) Name() string {
	return "test"
}

func (t tstFileInfo) Size() int64 {
	return int64(len(t.Msg))
}

func (t tstFileInfo) Mode() fs.FileMode {
	return 0666
}

func (t tstFileInfo) ModTime() time.Time {
	return time.Now()
}

func (t tstFileInfo) IsDir() bool {
	return t.Dir
}

func (t tstFileInfo) Sys() interface{} {
	return nil
}

func TestZipfile_Read(t *testing.T) {
	tests := []tstFileInfo{
		{
			Msg: []byte("Testing"),
			Dir: false,
		},
	}

	for _, tst := range tests {
		zf := zipfile{
			name: tst.Name(),
			info: tst,
			data: tst.Msg,
		}
		buf := make([]byte, len(tst.Msg))
		if _, err := zf.Read(buf); err != nil {
			t.Errorf("cannot read zipfile: %s", err)
		}
		if !bytes.Equal(buf, tst.Msg) {
			t.Errorf("read contents not equal: Actual(%s), Expected(%s)", string(buf), string(tst.Msg))
		}
	}
}

func TestZipfile_httpFS(t *testing.T) {
	mgr := NewManager(Options{ZipFileName: "./testdata/sqlfiles.zip"})
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(mgr)))

	t.Run("SqlFile", func(t *testing.T) {
		filePath := "/sqlfiles/default/selectDual.sql"
		req := httptest.NewRequest("GET", filePath, nil)
		wr := httptest.NewRecorder()

		mux.ServeHTTP(wr, req)

		resp := wr.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: Actual(%d), Expected(%d)", resp.StatusCode, http.StatusOK)
		}

		buf, err := os.ReadFile(path.Join("./testdata", filePath))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf, wr.Body.Bytes()) {
			t.Errorf("data not read correctly: Actual(%s), Expected(%s)", wr.Body.String(), string(buf))
		}
	})

	t.Run("Check directory", func(t *testing.T) {
		filePath := "/sqlfiles/"
		req := httptest.NewRequest("GET", filePath, nil)
		wr := httptest.NewRecorder()

		mux.ServeHTTP(wr, req)

		resp := wr.Result()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("unexpected status code: Actual(%d), Expected(%d)", resp.StatusCode, http.StatusOK)
		}
	})
}
