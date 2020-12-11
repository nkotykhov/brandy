package cpio

import (
	"time"
)

const (
	headerEOF = "TRAILER!!!"
)

type FileMode int64

type Header struct {
	DeviceID int
	Inode int64
	Mode FileMode
	UID int
	GID int
	Links int
	Mtime time.Time
	Size int64
	Name string
	NameSize int
	Linkname string
	Checksum uint32
}
