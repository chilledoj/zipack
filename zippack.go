package zipack

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"
)

// Manager is the container for the zipack logic
type Manager struct {
	cache       sync.Map
	zipFileName string
}

// interface check
var _ fs.FS = new(Manager)

// Options define the configuration for the manager object
type Options struct {
	ZipFileName  string
	PreloadPaths []string
}

// NewManager does exactly as it says on the tin.
func NewManager(opts Options) *Manager {
	mgr := &Manager{
		zipFileName: opts.ZipFileName,
	}
	if len(opts.PreloadPaths) > 1 {
		if err := mgr.preload(opts.PreloadPaths); err != nil {
			panic(err)
		}
	}
	return mgr
}

// GetFileContents will return the contents of the zipfile, using the cache first
func (mgr *Manager) GetFileContents(path string) ([]byte, error) {
	cacheFile, ok := mgr.cache.Load(path)
	if !ok {
		var err error
		cacheFile, err = mgr.readAndStoreFromZip(path)
		if err != nil {
			return nil, err
		}
	}

	file := cacheFile.(*zipfile)
	return file.data, nil
}

// GetReader will return an io.Reader for the zipfile, using the cache value first
func (mgr *Manager) GetReader(path string) (io.Reader, error) {
	val, ok := mgr.cache.Load(path)
	if ok {
		f, ok := val.(*zipfile)
		if !ok {
			return nil, fmt.Errorf("something went wrong")
		}
		return bytes.NewBuffer(f.data), nil
	}
	// Cache Miss - read from zip
	zf, err := mgr.readAndStoreFromZip(path)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(zf.data), nil
}

func (mgr *Manager) readAndStoreFromZip(path string) (*zipfile, error) {
	r, err := zip.OpenReader(mgr.zipFileName)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != path {
			continue
		}

		// let's not return anything and error for a directory
		fileInfo := f.FileHeader.FileInfo()
		if fileInfo.IsDir() {
			return nil, &os.PathError{
				Op:   "open",
				Path: path,
				Err:  fmt.Errorf("filepath is directory"),
			}
		}

		zipFile := &zipfile{
			name: path,
			info: fileInfo,
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer

		_, err = io.Copy(&buf, rc)
		if err != nil {
			return nil, err
		}
		rc.Close()

		zipFile.data = buf.Bytes()
		mgr.cache.Store(path, zipFile)
		return zipFile, nil
	}
	return nil, os.ErrNotExist
}

func (mgr *Manager) preload(paths []string) error {
	for _, path := range paths {
		_, err := mgr.readAndStoreFromZip(path)
		if err != nil {
			return err
		}
	}
	return nil
}
