package fingerprint

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var testDir string

func prepareTestDir() {
	var err error
	testDir, err = ioutil.TempDir("", "fingerprint-test")
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	prepareTestDir()
	defer (func() {
		if err := os.RemoveAll(testDir); err != nil {
			panic(err)
		}
	})()

	os.Exit(m.Run())
}

func TestCompile(t *testing.T) {
	testCases := []struct {
		content, fingerprint string
	}{
		{
			"a",
			"0cc175b9c0",
		},
		{
			"test",
			"098f6bcd46",
		},
	}

	for _, testCase := range testCases {
		var compiled = Compile([]byte(testCase.content))
		if compiled != testCase.fingerprint {
			t.Errorf(
				"compile failed for %s, got %s",
				testCase.content,
				compiled,
			)
		}
	}
}

func TestGetFileName(t *testing.T) {
	var testCases = []struct {
		path, fileName string
	}{
		{
			"/a/b/c.d",
			"c",
		},
		{
			"abc.d",
			"abc",
		},
		{
			"/this/one/no/ext",
			"ext",
		},
	}

	for _, testCase := range testCases {
		var fileName = getFileName(testCase.path)

		if fileName != testCase.fileName {
			t.Errorf(
				"get file name failed for %s, got %s",
				testCase.path,
				fileName,
			)
		}
	}
}

func createTestFile(t *testing.T, fileName, content string) *os.File {
	var (
		file *os.File
		err  error
	)

	file, err = os.Create(filepath.Join(testDir, fileName))
	if err != nil {
		t.Error(err)
	}
	if _, err = file.Write([]byte(content)); err != nil {
		t.Error(err)
	}
	if err = file.Close(); err != nil {
		t.Error(err)
	}

	file, err = os.Open(filepath.Join(testDir, fileName))
	if err != nil {
		t.Error(err)
	}

	return file
}

func TestFingerPrintedName(t *testing.T) {
	testCases := []struct {
		fileName, expectedFileName, content string
	}{
		{
			"a.js",
			"a-0cc175b9c0.js",
			"a",
		},
		{
			"test.css",
			"test-098f6bcd46.css",
			"test",
		},
	}

	for _, testCase := range testCases {
		f := createTestFile(t, testCase.fileName, testCase.content)
		file := makeFingerPrintedFile(f, "")

		name, err := file.FingerPrintedName()
		if err != nil {
			t.Error(err)
		}
		if name != testCase.expectedFileName {
			t.Errorf(
				"finger printed name failed for %s, got %s",
				testCase.expectedFileName,
				name,
			)
		}
	}
}

func TestFingerPrintedPath(t *testing.T) {
	testCases := []struct {
		fileName, expectedPath, destDir, content string
	}{
		{
			"a.js",
			filepath.Join(testDir, "a-0cc175b9c0.js"),
			"",
			"a",
		},
		{
			"test.css",
			filepath.Join("dest", "test-098f6bcd46.css"),
			"dest",
			"test",
		},
	}

	for _, testCase := range testCases {
		f := createTestFile(t, testCase.fileName, testCase.content)
		file := makeFingerPrintedFile(f, testCase.destDir)

		path, err := file.FingerPrintedPath()
		if err != nil {
			t.Error(err)
		}
		if path != testCase.expectedPath {
			t.Errorf(
				"finger printed path failed for %s, got %s",
				testCase.expectedPath,
				path,
			)
		}
	}
}

func prepareTestFiles(t *testing.T, contents []string) []string {
	testPaths := make([]string, len(contents))
	for i, content := range contents {
		testFile, err := ioutil.TempFile(testDir, "")
		if err != nil {
			t.Error(err)
		}
		if _, err = testFile.Write([]byte(content)); err != nil {
			t.Error(err)
		}

		testPaths[i] = testFile.Name()
	}

	return testPaths
}

func TestCompileAndWriteFiles(t *testing.T) {
	testPaths := prepareTestFiles(
		t,
		[]string{
			"test",
			"foobar",
		},
	)

	if err := CompileAndWriteFiles(testPaths, ""); err != nil {
		t.Error(err)
	}
}
