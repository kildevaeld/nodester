package nodester

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	NODE_REPO     = "https://nodejs.org/dist/"
	NODE_MANIFEST = NODE_REPO + "index.json"
)

type Config struct {
	Root    string
	Cache   string
	Source  string
	Current string
}

type RemoteOptions struct {
	Max                int
	Lts                bool
	HostCompatibleOnly bool
}

type Nodester struct {
	config Config
}

func (self *Nodester) List() ([]string, error) {

	path := self.config.Source

	if err := ensureDir(path); err != nil {
		return nil, err
	}

	return nil, nil
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

func (self *Nodester) Install() error {
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

func (self *Nodester) download(manifest Manifest, Os, arch string, source bool, progress func(progress, total int)) error {

	localFile := filepath.Join(self.config.Cache, manifest.localFile(os, arch, source))

	res, err := http.Get(manifest.remoteFile(Os, arch, source))

	if err != nil {
		return err
	}

	file, err := os.Open(localFile)

	if err != nil {
		return err
	}

	defer file.Close()
	defer res.Body.Close()

}

func (self *Nodester) Init() error {

	path := self.config.Root

	if !DirExists(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}

	conf := &self.config

	conf.Cache = filepath.Join(path, "cache")
	conf.Source = filepath.Join(path, "source")
	conf.Current = filepath.Join(path, "current")

	if err := ensureDir(conf.Cache); err != nil {
		return err
	}

	if err := ensureDir(conf.Source); err != nil {
		return err
	}

	if err := ensureDir(conf.Current); err != nil {
		return err
	}

	return nil
}

func New(config Config) *Nodester {
	return &Nodester{config}
}
