package cpio

import (
	"errors"
	"io"
)

type fileStreamLimitReader struct {
	r io.ReadSeeker
	start int64
	limit int64
	n int64
}


func newFileStreamLimitReader(r io.ReadSeeker, limit int64) (*fileStreamLimitReader, error) {
	start, err := r.Seek(0, 1)
	if err != nil {
		return nil, err
	}
	return &fileStreamLimitReader{
		r: r,
		limit: limit,
		start: start,
	}, nil
}

func (f *fileStreamLimitReader) Read(d []byte) (n int, err error) {
	if f.n >= f.limit {
		return 0, io.EOF
	}
	// out of order reads
	p, err := f.r.Seek(0, 1)
	if err != nil {
		return 0, err
	}
	if f.start+f.n != p {
		return 0, errors.New("error out of order read")
	}

	l := int64(len(d))
	if l > f.limit-f.n {
		d = d[0:f.limit-f.n]
	}
	n, err = f.r.Read(d)
	f.n+=int64(n)
	return
}
