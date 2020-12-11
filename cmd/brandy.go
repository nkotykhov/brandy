package main

import (
	"fmt"
	"github.com/rocky-linux/brandy/rpm"
	"github.com/rocky-linux/brandy/rpm/rpmutil"
	"io"
	"os"
	"path"
	"path/filepath"
)

func srpm2git(pkg *rpmutil.Package, basedir string) error {
	specDir := path.Join(basedir, "SPECS")
	if _, err := os.Stat(specDir); os.IsNotExist(err) {
		err := os.Mkdir(specDir, 0760)
		if err != nil {
			return err
		}
	}
	payload, err := pkg.Payload()
	if err != nil {
		return err
	}
	for {
		h, f, err := payload.Next()
		if err != nil {
			return err
		}
		if filepath.Ext(h.Name) == ".spec" {
			err := func() error {
				dest := path.Join(specDir, h.Name)
				out, err := os.Create(dest)
				if err != nil {
					return err
				}
				defer out.Close()
				_, err = io.Copy(out, f)
				if err != nil && err != io.EOF {
					panic(err)
					return err
				}
				return nil
			}()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage %s [file] [path]\n", os.Args[0])
		os.Exit(1)
	}

	if _, err := os.Stat(os.Args[2]); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "base directory doesn't exits\n")
		os.Exit(1)
	}

	pkgf, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open package file: %s\n", err.Error())
		os.Exit(1)
	}

	defer pkgf.Close()
	pkg, err := rpmutil.ReadPackage(pkgf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed during package parsing: %s\n", err.Error())
		os.Exit(1)
	}

	err = srpm2git(pkg, os.Args[2])
	if err != nil {
		panic(err)
	}
}
