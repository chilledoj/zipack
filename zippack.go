package zipack

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"sync"
)

// Manager is the container for the zipack logic
type Manager struct {
	cache       sync.Map
	zipFileName string
}

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

// GetFileContents will return the contents of the file, using the cache first
func (mgr *Manager) GetFileContents(path string) ([]byte, error) {
	_, err := mgr.GetReader(path)
	if err != nil {
		return nil, err
	}
	// use cache now
	v, _ := mgr.cache.Load(path)

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
	return mgr.readAndStoreFromZip(path)
}

func (mgr *Manager) readAndStoreFromZip(path string) (io.Reader, error) {
	r, err := zip.OpenReader(mgr.zipFileName)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var buf bytes.Buffer

	for _, f := range r.File {
		if f.Name != path {
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

func (mgr *Manager) preload(paths []string) error {
	for _, path := range paths {
		_, err := mgr.readAndStoreFromZip(path)
		if err != nil {
			return err
		}
	}
	return nil
}
