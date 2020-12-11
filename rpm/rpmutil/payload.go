package rpmutil

import (
	"github.com/rocky-linux/brandy/cpio"
)

type PayloadReader interface {
	Next() error
	Read(d []byte) (int, error)
}

type payloadReader struct {
	r cpio.Reader
}
