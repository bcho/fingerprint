// $ ./finger style.css
// $ ./finger .
// $ ./finger -r .
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bcho/fingerprint"
)

func main() {
	var (
		recursive   bool
		filesOrDirs []string
	)

	flag.BoolVar(&recursive, "r", false, "fingerprint directories recursively")
	flag.Parse()

	filesOrDirs = flag.Args()
	if len(filesOrDirs) < 1 {
		showUsage()
		os.Exit(0)
	}

	for _, src := range filesOrDirs {
		source, err := filepath.Abs(src)
		if err != nil {
			hanleError(err)
		}

		stat, err := os.Stat(source)
		if err != nil {
			hanleError(err)
		}

		if stat.IsDir() {
			if err = compileDir(source, recursive); err != nil {
				hanleError(err)
			}
		} else {
			if err = compileFile(source); err != nil {
				hanleError(err)
			}
		}
	}
}

func compileFile(source string) error {
	return fingerprint.CompileAndWriteFiles([]string{source})
}

func compileDir(source string, recursive bool) error {
	stats, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}

	var files, dirs []string
	for _, stat := range stats {
		fullPath := filepath.Join(source, stat.Name())

		if stat.IsDir() {
			if recursive {
				dirs = append(dirs, fullPath)
			}
		} else {
			// TODO check file type?
			files = append(files, fullPath)
		}
	}

	if err = fingerprint.CompileAndWriteFiles(files); err != nil {
		return err
	}

	for _, subPath := range dirs {
		if err = compileDir(subPath, true); err != nil {
			return err
		}
	}

	return nil
}

func showUsage() {
	fmt.Fprintf(os.Stderr, `Usage: finger [OPTIONS] SOURCE...
Fingerprinting files.

Options:
  -r                    fingerprinting directories recursively.
`)
}

// TODO handle error like a gentleman
func hanleError(err error) {
	panic(err)
}
