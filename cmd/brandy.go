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
		hdr, err := pr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Printf("fileName: %s\n", hdr.Name)
	}

	return nil

}

func main() {
	err := extractPackage(os.Args[1], "")
	if err != nil {
		panic(err)
	}
}
