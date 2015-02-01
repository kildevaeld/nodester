package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	NODE_MIRROR = "http://nodejs.org/dist/"
	IOJS_MIRROR = "https://iojs.org/dist/"
)

type NodeManager struct {
	path     string
	current  string
	platform string
	arch     string
}

func exists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

// Use a given node version
func (n *NodeManager) Use(version string) {
	version = normalizeVersion(version)

	if n.current == version {
		return
	}

	srcPath, _ := filepath.Abs(n.nodePath(&version))
	destPath, _ := filepath.Abs(filepath.Join(n.path, "current/"))

	if exists(destPath) {
		os.RemoveAll(destPath)
	}

	srcPath = srcPath + "/"

	os.Symlink(srcPath, destPath)

	cur := filepath.Join(n.path, "CURRENT_VERISON")

	n.current = version
	ioutil.WriteFile(cur, []byte(version), 0755)
}

func (n *NodeManager) Has(version string) bool {
	version = normalizeVersion(version)
	for _, v := range n.List() {
		if v == version {
			return true
		}
	}
	return false
}

func (n *NodeManager) Current() string {
	return n.current
}

func (n *NodeManager) List() []string {
	files, _ := ioutil.ReadDir(filepath.Join(n.path, "node"))

	var out []string
	for _, file := range files {
		if file.IsDir() {
			out = append(out, file.Name())
		}
	}
	return out
}

func (n *NodeManager) ListRemote() ([]string, error) {
	res, err := http.Get(NODE_MIRROR)
	if err != nil {
		return []string{}, err
	}
	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	version := string(b)

	r, _ := regexp.Compile("<a[\\s\\w=\".\\/]*>([v\\d.]*)\\/\\s*<\\/a>")

	v := r.FindAllStringSubmatch(version, -1)

	var versions []string

	for _, s := range v {
		if len(s) == 2 && s[1] != ".." {
			versions = append(versions, s[1])
		}
	}

	return versions, nil

}

func (n *NodeManager) Install(version string) error {
	version = normalizeVersion(version)

	if n.Has(version) {
		return nil
	}

	dest, err := n.Download(version, nil)

	if err != nil {
		return err
	}

	arch := n.arch
	if arch == "amd64" {
		arch = "x64"
	} else if arch == "386" {
		arch = "x86"
	}

	outputn := "node-" + version + "-" + n.platform + "-" + arch

	unpack_dest := filepath.Join(n.path, "node")

	cmd := exec.Command("tar", "-zxvf", dest, "-C", unpack_dest)

	cmd.Run()

	fp := filepath.Join(unpack_dest, version)
	os.Rename(filepath.Join(unpack_dest, outputn), fp)

	return nil

}

func (n *NodeManager) Remove(version string) error {
	version = normalizeVersion(version)
	if !n.Has(version) {
		return nil
	}

	dest := filepath.Join(n.path, "node", version)

	return os.RemoveAll(dest)

}

func (n *NodeManager) CleanCache() (err error) {
	dir := filepath.Join(n.path, "src")

	err = os.RemoveAll(dir)

	err = os.MkdirAll(dir, 0755)
	return err
}

func (n *NodeManager) Download(version string, fn func(progress DownloadProgress)) (string, error) {
	version = normalizeVersion(version)

	if n.Has(version) {
		return "", nil
	}

	arch := n.arch
	if arch == "amd64" {
		arch = "x64"
	} else if arch == "386" {
		arch = "x86"
	}

	outputn := "node-" + version + "-" + n.platform + "-" + arch
	filename := outputn + ".tar.gz"
	url := NODE_MIRROR + version + "/" + filename

	dest := filepath.Join(n.path, "src", filename)

	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		return dest, nil
	}

	out, _ := os.Create(dest)
	defer out.Close()

	_, e := DownloadSync(url, out, fn)

	return dest, e
}

func (n *NodeManager) init(path string) {
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
		os.Mkdir(filepath.Join(path, "node"), 0755)
		os.Mkdir(filepath.Join(path, "src"), 0755)
		//check(err)
	} else if !stat.IsDir() {
		return
	}

	n.path = path

	cur := filepath.Join(path, "CURRENT_VERISON")

	if _, err := os.Stat(cur); !os.IsNotExist(err) {
		str, _ := ioutil.ReadFile(cur)
		n.current = strings.Trim(string(str), " ")
	}

	n.platform = runtime.GOOS
	n.arch = runtime.GOARCH
}

func NewNodeManager(path string) *NodeManager {
	n := &NodeManager{}

	n.init(path)

	return n
}

func normalizeVersion(version string) string {
	version = strings.Trim(version, " ")
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	return version
}

func (n *NodeManager) sourcePath(version *string) string {
	path := filepath.Join(n.path, "src")

	if version == nil {
		return path
	}

	arch := n.arch
	if arch == "amd64" {
		arch = "x64"
	} else if arch == "386" {
		arch = "x86"
	}

	outputn := "node-" + *version + "-" + n.platform + "-" + arch
	filename := outputn + ".tar.gz"
	return filepath.Join(path, filename)
}

func (n *NodeManager) nodePath(version *string) string {
	str := filepath.Join(n.path, "node")
	if version == nil {
		return str
	}
	return filepath.Join(str, *version)
}
