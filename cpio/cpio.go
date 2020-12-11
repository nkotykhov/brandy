package cpio

import (
	"errors"
	"io"
	"io/ioutil"
)


type Reader interface {
	Next() (*Header,io.Reader,  error)
	Read([]byte) (int, error)
}

type streamReader struct {
	r *readSeekCounter
	next int64
}

func NewReader(r io.Reader) (Reader, error){
	return &streamReader{r: &readSeekCounter{r: r}}, nil
}


func (s *streamReader) Next() (*Header, io.Reader,  error) {
	if s.next != s.r.n {
		if _, err := s.r.Seek(s.next-s.r.n, 1); err != nil {
			return nil, nil, err
		}
	}
	h, err := ReadNewcHeader(s)
	if err != nil {
		return nil, nil, err
	}
	f, err := newFileStreamLimitReader(s.r, h.Size)
	if err != nil {
		return nil, nil, err
	}
	s.next = padding64(h.Size+s.r.n)

	return h, f, nil
}

func (r streamReader) Read(d []byte) (int, error) {
	return r.r.Read(d)
}


func padding(i int) int {
	return 3 + i - (i+3)%4
}

func padding64(i64 int64) int64 {
	return 3 + i64 - (i64+3)%4
}

type readSeekCounter struct {
	r io.Reader
	n int64
}

func (r *readSeekCounter) Read(d []byte) (n int, err error) {
	n, err = r.r.Read(d)
	r.n += int64(n)
	return
}

func (r *readSeekCounter) Seek(offset int64, whence int) (n int64, err error) {
	if whence != 1 {
		return 0, errors.New("error can only seek from current position")
	}
	if offset == 0 {
		return r.n, nil
	}

	return io.CopyN(ioutil.Discard, r, offset)
}
