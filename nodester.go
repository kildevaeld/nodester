package nodester

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	NODE_REPO     = "https://nodejs.org/dist/"
	NODE_MANIFEST = NODE_REPO + "index.json"
)

var (
	ErrInvalidVersion   = errors.New("Invalid version")
	ErrExistsNotInCache = errors.New("cache")
)

type Config struct {
	Root    string
	Cache   string
	Source  string
	Current string
	Temp    string
}

type RemoteOptions struct {
	Max                int
	Lts                bool
	HostCompatibleOnly bool
}

type Nodester struct {
	config Config
}

func (self *Nodester) List() (Manifests, error) {

	path := self.config.Source

	if err := ensureDir(path); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	var out Manifests
	for _, file := range files {

		if !file.IsDir() {
			continue
		}
		split := strings.Split(file.Name(), "-")

		if len(split) < 2 {
			continue
		}
		manifest, merr := self.GetManifest(split[1])
		if merr != nil {
			return nil, merr
		}
		manifest.Installed = true
		out = append(out, manifest)
	}

	return out, nil
}

func (self *Nodester) Current() string {
	path := filepath.Join(self.config.Root, "version")

	if FileExists(path) {
		bs, _ := ioutil.ReadFile(path)
		return string(bs)
	}
	return ""
}

func downloadManifests(config Config) (Manifests, error) {

	cache := filepath.Join(config.Cache, "manifests.json")

	if FileExists(cache) {
		s, _ := os.Stat(cache)
		fma := time.Now().Add(-5 * time.Minute)
		if !s.ModTime().Before(fma) {
			bs, err := ioutil.ReadFile(cache)

			if err == nil {
				var results Manifests
				if err := json.Unmarshal(bs, &results); err != nil {
					return nil, err
				}
				return results, err
			}
		}
	}

	res, err := http.Get(NODE_MANIFEST)

	if err != nil {
		return nil, err
	}
	var bs []byte
	bs, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	var results Manifests
	if err := json.Unmarshal(bs, &results); err != nil {
		return nil, err
	}

	ioutil.WriteFile(cache, bs, 0755)

	return results, nil
}

func (self *Nodester) ListRemote(options RemoteOptions) (Manifests, error) {
	var m Manifests
	var err error
	if m, err = downloadManifests(self.config); err != nil {
		return nil, err
	}

	max := len(m)
	ln := max

	if options.Max != 0 && options.Max <= max {
		max = options.Max
	}

	if options.HostCompatibleOnly {
		var found int
		for _, man := range m {
			if man.isHostCompatible(runtime.GOOS, runtime.GOARCH) {
				m[found] = man
				found++
			}
		}
		m = m[0:found]
	}

	for i, man := range m {
	inner:
		for _, v := range man.Files {

			split := strings.Split(v, "-")

			var arch, oss string

			if len(split) == 2 {
				arch = split[0]
				oss = split[1]
			} else {

			}

			if self.Has(Version{
				Version: man.Version,
				Arch:    arch,
				Os:      oss,
				Source:  len(split) == 1,
			}) {
				man.Installed = true
				break inner
			}
		}
		m[i] = man
	}

	if !options.Lts {
		if max == ln {
			return m, nil
		}
		return m[0:max], nil
	}

	var found = 0
	for _, man := range m {
		if man.isLts() {
			m[found] = man
			found++
		}
	}

	if max > found {
		max = found
	}

	return m[0:max], nil
}

func (self *Nodester) Remove(version Version) error {
	manifest, err := self.GetManifest(version.Version)

	if err != nil {
		return err
	}

	if manifest.Version == "" {
		return ErrInvalidVersion
	}

	sourceDir := filepath.Join(self.config.Source, version.Name())

	if manifest.Version == self.Current() {
		os.Remove(self.config.Current)
		os.Remove(filepath.Join(self.config.Root, "version"))
	}

	if DirExists(sourceDir) {
		return os.RemoveAll(sourceDir)
	}

	return nil
}

func (self *Nodester) ClearCache() error {
	os.RemoveAll(self.config.Cache)
	return os.MkdirAll(self.config.Cache, 0755)
}

func (self *Nodester) Install(version Version, progressCB func(step Step)) error {

	manifest, err := self.GetManifest(version.Version)

	if err != nil {
		return err
	}

	if manifest.Version == "" {
		return ErrInvalidVersion
	}

	if progressCB == nil {
		progressCB = func(Step) {}
	}

	sourceDir := filepath.Join(self.config.Source, version.Name())

	if !DirExists(sourceDir) {
		localFile := filepath.Join(self.config.Cache, manifest.localFile(version.Os, version.Arch, version.Source))

		if !FileExists(localFile) {
			return ErrExistsNotInCache
		}

		file, ferr := os.Open(localFile)

		if ferr != nil {
			return ferr
		}
		defer file.Close()

		reader, rerr := gzip.NewReader(file)
		if rerr != nil {
			return rerr
		}
		defer reader.Close()
		progressCB(Unpack)

		target := self.config.Source

		if version.Source {
			target = self.config.Temp
		}

		err = UnpackFile(reader, target, 0)

		if err != nil {
			return err
		}

		if version.Source {
			progressCB(Compile)
			return compile(filepath.Join(self.config.Temp, version.Name()), self.config, version, progressCB)
		}
	}

	return err
}

func (self *Nodester) Has(version Version) bool {
	return DirExists(filepath.Join(self.config.Source, version.Name()))
}

func (self *Nodester) Use(version Version) error {

	m, err := self.GetManifest(version.Version)

	if err != nil {
		return err
	}

	destPath := self.config.Current
	if destPath[len(destPath)-1] != '/' {
		//destPath += "/"
	}

	if DirExists(destPath) || FileExists(destPath) {
		if err := os.RemoveAll(destPath); err != nil {
			return err
		}
	}

	sourcePath := filepath.Join(self.config.Source, version.Name())

	if sourcePath[len(sourcePath)-1] != '/' {
		//sourcePath += "/"
	}

	if !DirExists(sourcePath) {
		return errors.New("Not installed")
	}

	err = os.Symlink(sourcePath, destPath)

	if err != nil {
		return err
	}

	versionFile := filepath.Join(self.config.Root, "version")

	ioutil.WriteFile(versionFile, []byte(m.Version), 0755)

	return nil
}

func (self *Nodester) Download(version Version, progressCB func(progress, total int64)) error {

	manifest, err := self.GetManifest(version.Version)

	if err != nil {
		return err
	}

	if manifest.Version == "" {
		return ErrInvalidVersion
	}

	localFile := filepath.Join(self.config.Cache, manifest.localFile(version.Os, version.Arch, version.Source))

	if !FileExists(localFile) {
		err = self.download(manifest, version.Os, version.Arch, version.Source, progressCB)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Nodester) GetManifest(version string) (Manifest, error) {
	version = normalizeVersion(version)

	manifests, err := downloadManifests(self.config)

	if err != nil {
		return Manifest{}, err
	}

	for _, man := range manifests {
		if man.Version == version {
			return man, nil
		}
	}
	return Manifest{}, nil
}

func (self *Nodester) download(manifest Manifest, Os, arch string, source bool, progressCB func(progress, total int64)) error {

	localFile := filepath.Join(self.config.Cache, manifest.localFile(Os, arch, source))

	log.Printf("Remote file %s\n", manifest.remoteFile(Os, arch, source))
	res, err := http.Get(manifest.remoteFile(Os, arch, source))

	if err != nil {
		return err
	}

	file, err := os.Create(localFile)

	if err != nil {
		return err
	}

	defer file.Close()
	defer res.Body.Close()

	_, err = copy(file, res.Body, func(progress int64) {
		if progressCB != nil {
			progressCB(progress, res.ContentLength)
		}
	})

	return err
}

func (self *Nodester) Init() error {
	var err error
	path := self.config.Root

	if !DirExists(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	path, err = filepath.Abs(path)

	if err != nil {
		return err
	}

	conf := &self.config

	conf.Cache = filepath.Join(path, "cache")
	conf.Source = filepath.Join(path, "source")
	conf.Current = filepath.Join(path, "current")
	conf.Temp = filepath.Join(path, "temp")
	if err := ensureDir(conf.Cache); err != nil {
		return err
	}

	if err := ensureDir(conf.Source); err != nil {
		return err
	}

	if err := ensureDir(conf.Current); err != nil {
		return err
	}

	if err := ensureDir(conf.Temp); err != nil {
		return err
	}

	return nil
}

func New(config Config) *Nodester {
	return &Nodester{config}
}
