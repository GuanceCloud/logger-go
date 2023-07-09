package logger

import "os"

type fileSyncer struct {
	f  string
	fd *os.File
}

func mustNewFileSyncer(f string) *fileSyncer {
	fs := &fileSyncer{
		f: f,
	}

	if fd, err := os.OpenFile(fs.f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		panic(err.Error())
	} else {
		fs.fd = fd
		return fs
	}
}

func (fs *fileSyncer) Write(p []byte) (int, error) {
	return fs.fd.Write(p)
}

func (fs *fileSyncer) Sync() error {
	return fs.fd.Close()
}
