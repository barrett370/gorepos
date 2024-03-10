package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"os"
	"path"
)

type Package struct {
	Repo     string
	Path     string
	Packages []string
	URL      string
}

func writePage(w io.Writer, path string, pkg Package) error {
	t, err := template.ParseFiles("template.html")
	if err != nil {
		return err
	}
	pkg.Path = path
	err = t.Execute(w, pkg)
	if err != nil {
		return err
	}
	return nil
}

func loadConfig(namespace string) (ret []Package, err error) {
	var f *os.File
	f, err = os.Open(fmt.Sprintf("%s.json", namespace))
	if err != nil {
		return
	}
	err = json.NewDecoder(f).Decode(&ret)
	return
}

func main() {
	flag.Parse()
	namespace := flag.Arg(0)
	root := flag.Arg(1)

	packages, err := loadConfig(namespace)
	if err != nil {
		log.Fatalf("%s\n", err)
		return
	}

	for _, pkg := range packages {
		packages := pkg.Packages
		packages = append(packages, "")
		for _, subpackage := range packages {
			pkgpath := path.Join(namespace, pkg.Path)

			root := path.Join(root, pkg.Path, subpackage)

			if err := os.MkdirAll(root, 0755); err != nil {
				slog.Error("error", "err", err)
				return
			}

			fd, err := os.Create(path.Join(root, "index.html"))
			if err != nil {
				slog.Error("error", "err", err)
				return
			}
			defer fd.Close()
			err = writePage(fd, pkgpath, pkg)
			if err != nil {
				slog.Error("error", "err", err)
			}
			fd.Close()
		}
	}
}

// vim: foldmethod=marker
