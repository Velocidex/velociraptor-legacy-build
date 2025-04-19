//go:build mage
// +build mage

package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/magefile/mage/sh"
)

type mutation struct {
	// Copy files
	From, To string

	// Go Replace by regex
	Match, Replace, Glob string

	DeleteGlob string
}

var (
	deps = []string{
		"https://github.com/Velocidex/WinPmem",
		"https://github.com/Velocidex/go-vhdx",
		"https://github.com/Velocidex/go-vmdk",
		"https://github.com/Velocidex/go-journalctl",
		"https://github.com/Velocidex/velociraptor",
	}

	golang_url = "https://go.dev/dl/go1.20.14.linux-amd64.tar.gz"
	mutations  = []mutation{
		{From: "../patches/go.mod", To: "velociraptor/go.mod"},
		{From: "../patches/go.sum", To: "velociraptor/go.sum"},
		{DeleteGlob: "velociraptor/tools/survey/*.go"},
		{Glob: "velociraptor/vql/psutils/*.go",
			Match:   "github.com/shirou/gopsutil/v4",
			Replace: "github.com/shirou/gopsutil/v3"},
		{From: "../patches/survey.go", To: "velociraptor/tools/survey/survey.go"},
		{From: "../patches/compat.go", To: "velociraptor/utils/compat.go"},
		{Glob: "*/go.mod", Match: "go 1.2", Replace: "// "},
		{Glob: "WinPmem/go-winpmem/go.mod", Match: "go 1.2", Replace: "// go 1.2"},
	}
)

func installGo() error {
	dst := "go"

	stat, err := os.Lstat(dst)
	if err == nil && stat.Mode().IsDir() {
		return nil
	}

	resp, err := http.Get(golang_url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil

		case err != nil:
			return err

		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			err := os.MkdirAll(target, 0755)
			if err != nil {
				return err
			}

		case tar.TypeReg:
			fmt.Printf("Creating %v (%v bytes)\n", target, header.Size)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			_, err = io.Copy(f, tr)
			if err != nil {
				return err
			}

			f.Close()

			err = os.Chmod(target, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
		}
	}
}

func replace_string_in_file(filename string, old string, new string) error {
	read, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	newContents := strings.Replace(string(read), old, new, -1)
	return ioutil.WriteFile(filename, []byte(newContents), 0644)
}

func maybeClone(dep string) error {
	base := filepath.Base(dep)
	_, err := os.Lstat(base)
	if err == nil {
		return nil
	}

	return sh.RunV("git", "clone", "--depth", "1", dep)
}

func copyOutput() error {
	basepath, pattern := doublestar.SplitPattern("./output/velociraptor*.exe")
	fsys := os.DirFS(basepath)
	matches, err := doublestar.Glob(fsys, pattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		filename := filepath.Join(basepath, match)
		ext := filepath.Ext(filename)
		bare := strings.TrimSuffix(filename, ext)
		base := filepath.Base(bare)
		dst := "../../output/" + base + "-legacy" + ext
		fmt.Printf("Replacing %v in %v\n", filename, dst)
		err := sh.Copy(dst, filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func Build() error {
	err := os.MkdirAll("build", 0700)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)

	err = os.Chdir("build")
	if err != nil {
		return err
	}

	err = installGo()
	if err != nil {
		return err
	}

	for _, dep := range deps {
		err = maybeClone(dep)
		if err != nil {
			return err
		}
	}

	for _, m := range mutations {
		if m.From != "" {
			fmt.Printf("Copying %v to %v\n", m.From, m.To)
			basedir := filepath.Dir(m.To)
			os.MkdirAll(basedir, 0755)

			err := sh.Copy(m.To, m.From)
			if err != nil {
				return err
			}
		}

		if m.DeleteGlob != "" {
			basepath, pattern := doublestar.SplitPattern(m.DeleteGlob)
			fsys := os.DirFS(basepath)
			matches, err := doublestar.Glob(fsys, pattern)
			if err != nil {
				return err
			}

			for _, match := range matches {
				filename := filepath.Join(basepath, match)
				fmt.Printf("Deleting %v in %v\n", m.Match, filename)
				err = os.Remove(filename)
				if err != nil {
					return err
				}
			}
		}

		if m.Glob != "" {
			basepath, pattern := doublestar.SplitPattern(m.Glob)
			fsys := os.DirFS(basepath)
			matches, err := doublestar.Glob(fsys, pattern)
			if err != nil {
				return err
			}

			for _, match := range matches {
				filename := filepath.Join(basepath, match)
				fmt.Printf("Replacing %v in %v\n", m.Match, filename)
				err = replace_string_in_file(filename, m.Match, m.Replace)
				if err != nil {
					return err
				}
			}
		}

	}

	// Build steps
	err = os.Chdir("velociraptor")
	if err != nil {
		return err
	}

	env := make(map[string]string)
	env["PATH"] = "../go/go/bin/:" + os.Getenv("PATH")
	env["GOPATH"] = ""

	go_path, err := filepath.Abs("../go/go/bin/go")
	if err != nil {
		return err
	}
	env["MAGEFILE_GOCMD"] = go_path
	env["MAGEFILE_VERBOSE"] = "1"

	err = sh.RunWithV(env, go_path, "version")
	if err != nil {
		return err
	}

	err = sh.RunWithV(env, go_path, "run", "-v", "./make.go", "windows")
	if err != nil {
		return err
	}

	err = sh.RunWithV(env, go_path, "run", "-v", "./make.go", "windowsx86")
	if err != nil {
		return err
	}

	err = copyOutput()
	if err != nil {
		return err
	}

	return nil
}
