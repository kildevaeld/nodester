package nodester

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

func normalizeOs(v string) string {
	switch v {
	case "osx":
		return "darwin"
	case "windows":
		return "win"
	default:
		return v
	}

}

func normalizeArch(v string) string {
	switch v {
	case "amd64":
		return "x64"
	case "386":
		return "x86"
	default:
		return v
	}

}

func compile(sourceDir string, config Config, version Version, stepCb func(step Step)) error {
	errOut, _ := os.Create(filepath.Join(config.Temp, version.Name()+"-build.error"))
	stdOut, _ := os.Create(filepath.Join(config.Temp, version.Name()+"-build.error"))

	defer errOut.Close()
	defer stdOut.Close()

	if stepCb == nil {
		stepCb = func(Step) {}
	}

	if runtime.GOOS == "windows" {

	} else {

		mkCmd := func(c string, args ...string) *exec.Cmd {
			cmd := exec.Command(c, args...)
			cmd.Dir = sourceDir
			cmd.Env = os.Environ()
			cmd.Stdout = stdOut
			cmd.Stderr = errOut
			return cmd
		}

		stepCb(Configure)
		target := filepath.Join(config.Source, version.Name())
		cmd := mkCmd("python", "./configure", "--prefix="+target)

		err := cmd.Run()
		if err != nil {
			return err
		}

		stepCb(Build)
		cmd = mkCmd("make")

		err = cmd.Run()

		if err != nil {
			return err
		}

		stepCb(Install)
		cmd = mkCmd("make", "install")

		err = cmd.Run()

		if err != nil {
			return err
		}

	}

	return nil

}
