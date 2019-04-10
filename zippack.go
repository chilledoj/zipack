package zipack

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Manager is the container for the zipack logic
type Manager struct {
	cache        sync.Map
	zipFileName  string
	PreloadPaths []string // TODO
}

// NewManager does exactly as it says on the tin.
func NewManager(zipFileName string) *Manager {
	mgr := &Manager{
		zipFileName: zipFileName,
	}
	return mgr
}

// GetFileContents will return the contents of the file, using the cache first
func (mgr *Manager) GetFileContents(path string) ([]byte, error) {
	_, err := mgr.GetReader(path)
	if err != nil {
		return nil, err
	}
	// use cache now
	v, ok := mgr.cache.Load(path)
	if !ok {
		// Should never happen
		return nil, fmt.Errorf("cache not populated after getReader")
	}

	return v.([]byte), nil
}

// GetReader will return an io.Reader for the file, using the cache value first
func (mgr *Manager) GetReader(path string) (io.Reader, error) {
	val, ok := mgr.cache.Load(path)
	if ok {
		reader := bytes.NewBuffer(val.([]byte))
		return reader, nil
	}
	// Cache Miss - read from zip
	r, err := zip.OpenReader(mgr.zipFileName)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var buf bytes.Buffer

	// This is just to make sure we get the correct path.
	// Using the following structure
	//  + sqlfiles
	//    + default
	//      selectDual.sql
	//
	// zipping the sqlfiles folder means that the folder becomes part
	// of the path within the zip file.
	//
	// Thus we need to get the zipfile name and use as a folder name.
	//
	// e.g.
	//  default/selectDual.sql in sqlfiles.zip   =   sqlfiles/default/selectDual.sql
	//
	zipfile := filepath.Base(mgr.zipFileName)
	zipext := filepath.Ext(mgr.zipFileName)
	basePath := zipfile[0 : len(zipfile)-len(zipext)]
	filename := filepath.Join(basePath, path)

	for _, f := range r.File {
		if f.Name != filename {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(&buf, rc)
		if err != nil {
			return nil, err
		}
		rc.Close()
		mgr.cache.Store(path, buf.Bytes())
		return &buf, nil
	}
	return nil, os.ErrNotExist
}
