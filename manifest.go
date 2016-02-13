package nodester

import "fmt"

type Manifests []Manifest

type Manifest struct {
	Version   string
	Date      string
	Files     []string
	Lts       interface{}
	Modules   string
	Npm       string
	V8        string
	Uv        string
	Zlib      string
	OpenSSL   string
	Installed bool
}

func (self Manifest) isLts() bool {
	switch n := self.Lts.(type) {
	case bool:
		return n
	case string:

		if n == "false" || n == "False" {
			return false
		}
		return true
	}
	return false
}

func (self Manifest) isHostCompatible(oss, arch string) bool {

	oss = normalizeOs(oss)
	arch = normalizeArch(arch)
	if oss == "darwin" {
		oss = "osx"
	}

	ossarch := oss + "-" + arch
	if oss == "osx" {
		ossarch += "-tar"
	}

	for _, src := range self.Files {

		if ossarch == src {
			return true
		}
	}
	return false
}

func (self Manifest) remoteFile(oss, arch string, source bool) string {
	oss = normalizeOs(oss)
	arch = normalizeArch(arch)
	if arch == "win" && oss == "x64" {

	}

	if oss == "osx" {
		oss = "darwin"
	}

	fn := fmt.Sprintf("%s/%s/node-%s", NODE_REPO, self.Version, self.Version)
	if source {
		return fn + ".tar.gz"
	}

	return fmt.Sprintf("%s-%s-%s.tar.gz", fn, oss, arch)
}

func (self Manifest) checksumFile() string {

	return fmt.Sprintf("%s/%s/SHASUMS256.txt", NODE_REPO, self.Version)

}

func (self Manifest) localFile(oss, arch string, source bool) string {
	oss = normalizeOs(oss)
	arch = normalizeArch(arch)
	if arch == "win" && oss == "x64" {

	}

	if oss == "osx" {
		oss = "darwin"
	}

	fn := "node-" + self.Version
	if source {
		return fn + ".tar.gz"
	}

	return fmt.Sprintf("%s-%s-%s.tar.gz", fn, oss, arch)
}

type Version struct {
	Version string
	Arch    string
	Os      string
	Source  bool
}

func (self Version) Name() string {
	oss := normalizeOs(self.Os)
	arch := normalizeArch(self.Arch)
	if arch == "win" && oss == "x64" {

	}
	fn := "node-" + normalizeVersion(self.Version)
	if self.Source {
		return fn
	}

	return fmt.Sprintf("%s-%s-%s", fn, oss, arch)
}
