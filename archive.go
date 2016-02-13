package nodester

import (
	"path/filepath"

	"archive/tar"

	"fmt"
	"io"
	"os"
)

// Copy tarball `reader` to `path`.
func UnpackFile(reader io.Reader, path string, strip int) error {

	tarball := tar.NewReader(reader)

	for {
		header, err := tarball.Next()

		if err == io.EOF {

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
				return err
			}

		case tar.TypeReg:
			// handle normal file

			writer, err := os.Create(filename)

			if err != nil {
				return err
			}

			io.Copy(writer, tarball)
			writer.Close()
			err = os.Chmod(filename, os.FileMode(header.Mode))

			if err != nil {
				return err
			}

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
