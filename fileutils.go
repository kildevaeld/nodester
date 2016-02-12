package nodester

import (
	"io"
	"os"
)

func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

func DirExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func ensureDir(path string) error {
	if !DirExists(path) {
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

func copy(dest io.Writer, src io.Reader, progress func(p int64)) (written int64, err error) {

	buf := make([]byte, 32*1024)

	for {
		nr, er := src.Read(buf)

		if nr > 0 {
			nw, ew := dest.Write(buf[0:nr])

			if nw > 0 {
				written += int64(nw)
				progress(written)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}

		if er != nil {
			err = er
			break
		}

	}

	return written, err
}

func normalizeVersion(v string) string {
	if v[0] == 'v' {
		return v
	}
	return "v" + v
}
