package cpio

import (
	"io"
	"io/ioutil"
)


type Reader interface {
	Next() (*Header, error)
	Read([]byte) (int, error)
}

type reader struct {
	r *readCounter
	next int64
}

func NewReader(r io.Reader) (Reader, error){
	return &reader{r: &readCounter{r: r}}, nil
}


func (r *reader) Next() (*Header, error) {
	if r.next != r.r.n {
		_, err := io.CopyN(ioutil.Discard, r, r.next-r.r.n)
		if err != nil {
			return nil, err
		}
	}
	h, err := ReadNewcHeader(r)
	if err != nil {
		return nil, err
	}
	r.next = padding64(h.Size+r.r.n)
	return h, nil
}

func (r reader) Read(d []byte) (int, error) {
	return r.r.Read(d)
}

func padding(i int) int {
	return 3 + i - (i+3)%4
}

func padding64(i64 int64) int64 {
	return 3 + i64 - (i64+3)%4
}

type readCounter struct {
	r io.Reader
	n int64
}

func (r *readCounter) Read(d []byte) (n int, err error) {
	n, err = r.r.Read(d)
	r.n += int64(n)
	return
}
