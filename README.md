# fingerprint

[![Build Status](https://travis-ci.org/bcho/fingerprint.svg?branch=master)](https://travis-ci.org/bcho/fingerprint)

A simple library for fingerprinting files (with MD5).


## Usage

```go
import "github.com/bcho/fingerprint"

var compiled, _ := fingerprint.CompileFiles([]string{
	"/path/to/assets/css/style.css",
	"/path/to/assets/javascript/app.js",
}, "")

for _, file := range compiled {
	print(file.FingerPrintedPath())
}

fingerprint.CompileAndWriteFiles([]string{
	"/path/to/assets/css/style.css",
	"/path/to/assets/javascript/app.js",
}, "")
```

## API

- `Compile`
- `CompileFiles`
- `CompileAndWriteFiles`


## Commandline

### Build

```shell
$ git clone https://github.com/bcho/fingerprint.git
$ go build fingerprint/cmd/finger.go
```


### Usage

Check out:

```shell
$ ./finger
```


## License

[SMPPL](LICENSE)
