// $ ./finger style.css
// $ ./finger .
// $ ./finger -r .
// $ ./finger -r . -o dest/
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
		destDir     string
		filesOrDirs []string
	)

	flag.BoolVar(&recursive, "r", false, "")
	flag.StringVar(&destDir, "o", "", "")
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
			if err = compileDir(source, destDir, recursive); err != nil {
				hanleError(err)
			}
		} else {
			if err = compileFile(source, destDir); err != nil {
				hanleError(err)
			}
		}
	}
}

func compileFile(source, destDir string) error {
	return fingerprint.CompileAndWriteFiles([]string{source}, destDir)
}

func compileDir(source, destDir string, recursive bool) error {
	stats, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}

	var files []string
	var dirs []struct{ src, dest string }
	for _, stat := range stats {
		name := stat.Name()
		fullPath := filepath.Join(source, name)

		if stat.IsDir() {
			if recursive {
				dir := struct{ src, dest string }{
					fullPath,
					filepath.Join(destDir, name),
				}
				dirs = append(dirs, dir)
			}
		} else {
			// XXX check file type?
			files = append(files, fullPath)
		}
	}

	if err = fingerprint.CompileAndWriteFiles(files, destDir); err != nil {
		return err
	}

	for _, subPath := range dirs {
		if err = compileDir(subPath.src, subPath.dest, true); err != nil {
			return err
		}
	}

	return nil
}

func showUsage() {
	fmt.Fprintf(os.Stderr, `Usage: finger [OPTIONS] SOURCE...
Fingerprinting files.

Arguments & options:
  -r                    fingerprinting directories recursively.
  -o DEST               output directory, defaults to original directory.
`)
}

// TODO handle error like a gentleman
func hanleError(err error) {
	panic(err)
}
