package zipack

import (
	"io"
	"io/fs"
	"os"
)

func (mgr *Manager) Open(filename string) (fs.File, error) {

	zipfile, err := mgr.readAndStoreFromZip(filename)
	if err != nil {
		if verr, ok := err.(*os.PathError); ok {
			return nil, verr
		}
		return nil, err
	}
	return zipfile, nil
}

type zipfile struct {
	name   string
	info   fs.FileInfo
	data   []byte
	offset int64
}

// interface check
var _ fs.File = new(zipfile)

func (z zipfile) Stat() (fs.FileInfo, error) {
	return z.info, nil
}

func (z *zipfile) Read(p []byte) (n int, err error) {
	if z.offset >= int64(len(z.data)) {
		err = io.EOF
		return
	}
	n = copy(p, z.data[z.offset:])
	z.offset += int64(n)
	return
}

func (z *zipfile) Close() error {
	z.offset = 0
	return nil
}
