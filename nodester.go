package nodester

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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
	return nil, nil
}

func getManifest() (Manifests, error) {
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

	return results, nil
}

func (self *Nodester) ListRemote(options RemoteOptions) (Manifests, error) {
	var m Manifests
	var err error
	if m, err = getManifest(); err != nil {
		return nil, err
	}

	max := len(m)
	ln := max

	if options.Max != 0 || options.Max <= max {
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

	conf := self.config

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
