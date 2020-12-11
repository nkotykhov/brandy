package main

import (
	"fmt"
	"github.com/rocky-linux/brandy/rpm"
	"github.com/rocky-linux/brandy/rpm/rpmutil"
	"os"
)

// extractPackage extracts contents of RPM to dir
// using rpm2cpio and cpio utils
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
	return nil

}

func main() {
	err := extractPackage(os.Args[1], "")
	if err != nil {
		panic(err)
	}
}
