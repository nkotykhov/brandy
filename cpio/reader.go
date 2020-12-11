package cpio

import (
	"errors"
	"io"
	"io/ioutil"
	"strconv"
	"time"
)

const (
	newcMagic = "070701"
	newcHeaderLen = 110
)

type newcReader struct {
	r io.Reader
	sink [8]byte
}


func (r newcReader) Read(b []byte) (int, error) {
	return r.r.Read(b)
}

func (r newcReader) Read16() (int, error) {
	buf := r.sink[:8]
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	i16, err := strconv.ParseInt(string(buf), 16, 0)
	if err != nil {
		return 0, err
	}
	return int(i16), nil
}


func ReadNewcHeader(r io.Reader) (*Header, error) {
	cr := newcReader{r: r}
	magic := make([]byte, 6)
	if _, err := io.ReadFull(cr, magic);err != nil {
		return nil, err
	}

	if string(magic) != newcMagic {
		return nil, errors.New("Bad Header Magic")
	}

	h := &Header{}

	i,err := cr.Read16()
	if err != nil {
		return nil, err
	}
	h.Inode = int64(i)

	i, err = cr.Read16()
	if err != nil {
		return nil, err
	}
	h.Mode = FileMode(i)

	h.UID, err = cr.Read16()
	if err != nil {
		return nil, err
	}

	h.GID, err = cr.Read16()
	if err != nil {
		return nil, err
	}

	h.Links, err = cr.Read16()
	if err != nil {
		return nil, err
	}

	i, err = cr.Read16()
	if err != nil {
		return nil, err
	}
	h.Mtime = time.Unix(int64(i), 0)

	i, err = cr.Read16()
	if err != nil {
		return nil, err
	}
	h.Size = int64(i)

	// skip dev and rdev
	if _, err = io.CopyN(ioutil.Discard, cr, int64(32)); err != nil {
		return nil, err
	}

	i, err = cr.Read16()
	if err != nil {
		return nil, err
	}
	h.NameSize = i

	i, err = cr.Read16()
	if err != nil {
		return nil, err
	}
	h.Checksum = uint32(i)
	name := make([]byte, h.NameSize)
	if _, err = io.ReadFull(cr, name); err != nil {
		return nil, err
	}

	h.Name = string(name[:h.NameSize-1])
	if h.Name == headerEOF {
		return nil, io.EOF
	}
	pad := padding(newcHeaderLen+h.NameSize) - h.NameSize - newcHeaderLen
	if pad >0 {
		if _, err := io.CopyN(ioutil.Discard, cr, int64(pad)); err != nil {
			return nil, err
		}
	}
	return h, nil
}

