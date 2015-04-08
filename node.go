package main

import (
	"compress/gzip"
	"fmt"
	. "github.com/visionmedia/go-unpack"
	"io/ioutil"
	"os"
	_ "os/exec"
	"path/filepath"
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
	Services Service
}

//
func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Use a given node version
func (n *NodeManager) Use(version string) {

	v, p := n.Services.GetPrefix(version)

	if n.current == v {
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

	n.current = p + "@" + v
	ioutil.WriteFile(cur, []byte(n.current), 0755)
}

func (n *NodeManager) Has(version string) bool {
	vs, p := n.Services.GetPrefix(version)
	vap := p + "@" + vs
	for _, v := range n.List() {
		if v == vap {
			return true
		}
	}
	return false
}

func (n *NodeManager) Current() string {
	return n.current
}

func (n *NodeManager) List() []string {
	files, _ := ioutil.ReadDir(n.nodePath(nil))

	var out []string
	for _, file := range files {
		if file.IsDir() {
			out = append(out, strings.Replace(file.Name(), "-", "@", 1))
		}
	}
	return out
}

func (n *NodeManager) ListRemote() (versions string, err error) {
	return n.Services.ListRemote()
}

func (n *NodeManager) Install(version string) error {

	if n.Has(version) {
		return nil
	}

	dest, err := n.Download(version, nil)

	if err != nil {
		return err
	}

	unpack_p := filepath.Base(dest)
	unpack_p = strings.Replace(unpack_p, ".tar.gz", "", 1)

	v, p := n.Services.GetPrefix(version)

	dest_p := n.nodePath(nil)

	file, _ := os.Open(dest)
	reader, _ := gzip.NewReader(file)
	UnpackTarball(reader, dest_p, 0)

	unpack_p = filepath.Join(dest_p, unpack_p)
	rename_p := filepath.Join(dest_p, fmt.Sprintf("%s-%s", p, v))

	os.Rename(unpack_p, rename_p)

	return nil

}

func (n *NodeManager) Remove(version string) error {

	if !n.Has(version) {
		return nil
	}

	dest := n.nodePath(&version)
	return os.RemoveAll(dest)

}

func (n *NodeManager) CleanCache() (err error) {
	dir := n.sourcePath(nil)
	err = os.RemoveAll(dir)
	err = os.MkdirAll(dir, 0755)

	return err
}

func (n *NodeManager) Download(version string, fn func(progress DownloadProgress)) (string, error) {

	if n.Has(version) {
		return "", nil
	}

	url, filename := n.Services.RemoteFile(version, n.arch, n.platform)

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
	n.arch = normalizeArch(runtime.GOARCH)
}

func (n *NodeManager) normalizeVersion(version string) string {
	//s, _ := n.Services.NormalizeVersion(version)
	return ""
}

func NewNodeManager(path string) *NodeManager {
	n := &NodeManager{}
	n.Services = &NodeService{}
	n.init(path)

	return n
}

func normalizeArch(arch string) (a string) {
	if arch == "amd64" {
		a = "x64"
	} else if arch == "386" {
		a = "x86"
	}
	return
}

func (n *NodeManager) sourcePath(version *string) string {
	path := filepath.Join(n.path, "src")

	if version == nil {
		return path
	}

	arch := normalizeArch(n.arch)

	outputn := "node-" + *version + "-" + n.platform + "-" + arch
	filename := outputn + ".tar.gz"
	return filepath.Join(path, filename)
}

func (n *NodeManager) nodePath(version *string) string {

	str := filepath.Join(n.path, "node")
	if version == nil {
		return str
	}

	v, p := n.Services.GetPrefix(*version)

	return filepath.Join(str, p+"-"+v)
}

func (n *NodeManager) servicePath(service Service, version *string) string {
	str := filepath.Join(n.path, service.Name())
	return str
}
