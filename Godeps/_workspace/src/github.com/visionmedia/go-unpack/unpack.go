package unpack

import . "github.com/visionmedia/go-debug"
import "path/filepath"
import "archive/tar"
import "strings"
import "bufio"
import "io"
import "os"

var debug = Debug("unpack")

// Copy tarball `reader` to `path`.
func UnpackTarball(reader io.ReadCloser, path string, strip int) error {
	debug("unpacking to '%s'", path)
	tarball := tar.NewReader(reader)

	for {
		header, err := tarball.Next()

		if err == io.EOF {
			debug("eof")
			break
		}

		if err != nil {
			return err
		}

		info := header.FileInfo()

		dir := filepath.Dir(header.Name)
		dir = strings.Join(strings.Split(dir, "/")[strip:], "/")
		dir = filepath.Join(path, dir)
		dst := filepath.Join(dir, info.Name())

		if info.IsDir() {
			debug("mkdir '%s'", dir)
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
			continue
		}

		debug("unpack '%s' -> '%s'", header.Name, dst)
		file, err := os.Create(dst)
		if err != nil {
			return err
		}

		_, err = io.Copy(bufio.NewWriter(file), tarball)
		if err != nil {
			return err
		}
	}

	return nil
}
