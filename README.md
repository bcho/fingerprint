# fingerprint

A simple library for fingerprinting files.


## Usage

```go
import "github.com/bcho/fingerprint"

var compiled, err := fingerprint.CompileFiles([]string{
	"/path/to/assets/css/style.css",
	"/path/to/assets/javascript/app.js",
})

// compiled =>
// {
//      "/path/to/assets/css/style.css": {
//          "fingerprint": "deadbeef1234",
//          "content": "p { margin: 0; }"
//      },
//      "/path/to/assets/javascript/app.js": {
//          "fingerprint": "bba3e2",
//          "content": ";(function() { console.log('hello, world'); })();"
//      }
// }
```

## API

- `Compile`
- `CompileFiles`
- `CompileAndWriteFiles`
