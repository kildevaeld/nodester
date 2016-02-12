package nodester

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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
	Max int
	Lts bool
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

func getManifest(config Config) (Manifests, error) {

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
	if m, err = getManifest(self.config); err != nil {
		return nil, err
	}

	max := len(m)
	ln := max

	if options.Max != 0 && options.Max <= max {
		max = options.Max
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
