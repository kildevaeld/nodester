package main

import (
	"path/filepath"

	. "github.com/visionmedia/go-debug"
)

import "archive/tar"

import "io"
import "os"
import "fmt"

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

		filename := header.Name
		filename = filepath.Join(path, filename)
		filename, _ = filepath.Abs(filename)

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory

			err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		case tar.TypeReg:
			// handle normal file

			writer, err := os.Create(filename)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			io.Copy(writer, tarball)

			err = os.Chmod(filename, os.FileMode(header.Mode))

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			writer.Close()
		case tar.TypeSymlink:

			dirname := filepath.Dir(filename)

			basename := filepath.Base(filename)

			cur, err := os.Getwd()

			if err != nil {
				return err
			}

			os.Chdir(dirname)

			err = os.Symlink(header.Linkname, basename)
			if err != nil {
				return err
			}

			os.Chdir(cur)

		default:
			return fmt.Errorf("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
	}

	return nil
}
