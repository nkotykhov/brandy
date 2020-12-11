package rpmutil

import (
	"fmt"
	"compress/gzip"
	"compress/bzip2"
	"io"
	"github.com/rocky-linux/brandy/rpm"
)

const (
	plUncompressed = "uncompressed"
	plGzip = "gzip"
	plBzip2 = "bzip2"
)

func decompressPkgPayload(p *Package) (io.Reader, error) {
	var compressor string
	t, val, err := p.Header.GetTag(rpm.TagPayloadCompressor)

	if err != nil && t != rpm.DataTypeNotFound {
		return nil, err
	}
	compressor = plGzip
	if t != rpm.DataTypeNotFound {
		compressor = string(val)
	}

	switch compressor {
		case plGzip:
			return gzip.NewReader(p.r)
		case plBzip2:
			return bzip2.NewReader(p.r), nil
	}
	return nil, fmt.Errorf("unsuppored compression format: %s", compressor)
}
