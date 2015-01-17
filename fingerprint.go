// A simple library for fingerprinting files.
package fingerprint

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	length           = 10             // Fingerpriting length.
	magicFingerPrint = "feeling_good" // Initial fingerprinting
	magicDestDir     = ""             // Use original parent dir by setting destDir to empty string.
)

var hasher = md5.Sum

// Hashed file.
type FingerPrintedFile struct {
	fingerPrint, destDir string
	content              []byte
	*os.File
}

// Generate fingerprint for something.
func Compile(content []byte) string {
	var hashed = hasher(content)

	return hex.EncodeToString(hashed[:])[:length]
}

// Generate fingerprints for files.
func CompileFiles(filePaths []string, destDir string) ([]*FingerPrintedFile, error) {
	var files []*FingerPrintedFile

	for _, path := range filePaths {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		files = append(files, makeFingerPrintedFile(f, destDir))
	}

	return files, nil
}

// Generate fingerprints and write it to files.
//
// If destDir is empty string, use the same dir as the original file.
func CompileAndWriteFiles(files []string, destDir string) error {
	hashedFiles, err := CompileFiles(files, destDir)
	if err != nil {
		return err
	}

	for _, file := range hashedFiles {
		info, err := file.Stat()
		if err != nil {
			return err
		}

		path, err := file.FingerPrintedPath()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, file.content, info.Mode())
		if err != nil {
			return err
		}
	}

	return nil
}

func makeFingerPrintedFile(f *os.File, destDir string) *FingerPrintedFile {
	return &FingerPrintedFile{
		magicFingerPrint,
		destDir,
		nil,
		f,
	}
}

func (file *FingerPrintedFile) shouldCompiled() error {
	if file.fingerPrint == magicFingerPrint {
		content, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		file.content = content
		file.fingerPrint = Compile(content)
	}

	return nil
}

func (file *FingerPrintedFile) FingerPrintedName() (string, error) {
	if err := file.shouldCompiled(); err != nil {
		return "", err
	}

	_, fileNameWithExt := filepath.Split(file.Name())
	fileExt := filepath.Ext(fileNameWithExt)
	fileName := getFileName(fileNameWithExt)

	return fmt.Sprintf(
		"%s-%s%s", // TODO custom fmt
		fileName,
		file.fingerPrint,
		fileExt,
	), nil
}

func (file *FingerPrintedFile) FingerPrintedPath() (string, error) {
	parentDir := file.destDir
	if parentDir == magicDestDir {
		parentDir, _ = filepath.Split(file.Name())
	}

	fileName, err := file.FingerPrintedName()
	if err != nil {
		return "", err
	}

	return filepath.Join(parentDir, fileName), nil
}

// Get filename from file path.
func getFileName(path string) string {
	var (
		fileName = filepath.Base(path)
		ext      = filepath.Ext(path)
	)

	if ext != "" {
		fileName = strings.Split(fileName, ext)[0]
	}

	return fileName
}
