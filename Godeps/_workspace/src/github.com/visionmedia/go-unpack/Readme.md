
# go-unpack

 Tarball unpacking utility.

## Installation

```
$ go get github.com/visionmedia/go-unpack
```

## Example

 Example which copies some.tgz to "some/path",
 creating the destination directories if required,
 stripping 0 components.

```go
import . "github.com/visionmedia/go-unpack"
import "compress/gzip"
import "os"

func main() {
  file, _ := os.Open("some.tgz")
  reader, _ := gzip.NewReader(file)
  err := UnpackTarball(reader, "some/path", 0)
}
```

# License

 MIT