package main

import (
	"fmt"
	"io"
	"github.com/rocky-linux/brandy/rpm"
	"github.com/rocky-linux/brandy/rpm/rpmutil"
	"os"
)

func extractPackage(pkg, dir string) error {
	f, err := os.Open(pkg)
	if err != nil {
		return err
	}
	defer f.Close()
	p, err := rpmutil.ReadPackage(f)

	if err != nil {
		return err
	}
	_, name, err := p.Header.GetTag(rpm.TagName)
	if err != nil {
		return err
	}
	_, version, err := p.Header.GetTag(rpm.TagVersion)
	fmt.Printf("%s-%s\n", name, version)
	pr, err := p.Payload();

	if err != nil {
		return err
	}
	for {
		hdr,f, err := pr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Printf("fileName: %s\n", hdr.Name)
		fheader := make([]byte, 8);
		if _, err := f.Read(fheader); err != nil {
			fmt.Printf("unable to read file: %s\n", err.Error())
			continue
		}
		fmt.Printf("first 8 bytes are: %s\n", fheader)
	}

	return nil

}

func main() {
	err := extractPackage(os.Args[1], "")
	if err != nil {
		panic(err)
	}
}
