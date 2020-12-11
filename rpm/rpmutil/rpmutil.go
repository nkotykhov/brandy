package rpmutil

import (
	"errors"
	"encoding/binary"
	"io"
	"io/ioutil"
	"github.com/rocky-linux/brandy/rpm"
	"github.com/rocky-linux/brandy/cpio"
)

type Package struct {
	SigHeader *rpm.Header
	Header *rpm.Header
	r *readCounter
}

func ReadPackage(r io.Reader) (*Package, error) {
	rc := &readCounter{r: r}
	// first lets get rid of Lead, for sanity check
	// will still use first LeadMagic to make sure
	// file is considered to be RPM, but rpm itself
	// no longe uses lead structure. it's there for
	// `file` command

	lead := make([]byte, rpm.LeadSize)
	if _, err:= io.ReadFull(rc, lead); err != nil {
		return nil, err
	}

	magic := binary.BigEndian.Uint32(lead[0:4])
	if magic&0xFFFFFFFF != rpm.LeadMagic {
		return nil,  errors.New("bad lead magic")
	}

	sigHeader, err := rpm.ReadHeader(rc)
	if err != nil {
		return nil, err
	}

	// signature header padded to align to 8 bytes
	psize := (rc.n +7) / 8 * 8
	skip := int64(psize-rc.n)

	if _, err := io.CopyN(ioutil.Discard, rc, skip); err != nil {
		return nil, err
	}

	header, err := rpm.ReadHeader(rc);

	if err != nil {
		return nil, err
	}
	pkg := &Package{
		SigHeader: sigHeader,
		Header: header,
		r: rc,
	}
	return pkg, nil
}

func (pkg *Package) Payload() (cpio.Reader, error) {
	plRdr, err := decompressPkgPayload(pkg)
	if err != nil {
		return nil, err
	}
	return cpio.NewReader(plRdr)
}


func (pkg *Package) Files() ([]FileInfo, error) {

	paths, err := pkg.Header.GetStrings(rpm.TagDirNames)
	if err != nil {
		return nil, err
	}

	files := make([]FileInfo, len(paths))
	for i:=0;i<len(files);i++ {
		files[i].Name = paths[i]
	}
	return files, nil
}


type readCounter struct {
	n int
	r io.Reader
}


func (rc *readCounter) Read(b []byte) (n int,err error) {
	n, err = rc.r.Read(b)
	rc.n += n
	return
}
