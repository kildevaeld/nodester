package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type Service interface {
	Name() string
	ListRemote() (string, error)
	RemoteFile(version string, arch string, platform string) (string, string)
	GetPrefix(str string) (version string, prefix string)
}

type NodeService struct{}

func (n *NodeService) Name() string { return "node" }

func (n *NodeService) ListRemote() (string, error) {
	var lock sync.Mutex
	result := make(map[string]interface{})

	var wg sync.WaitGroup

	urls := []string{NODE_MIRROR, IOJS_MIRROR}

	for _, url := range urls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			res, err := http.Get(url)
			defer res.Body.Close()

			lock.Lock()
			defer lock.Unlock()

			if err != nil {
				result[url] = err
			} else {
				b, e := ioutil.ReadAll(res.Body)

				if e != nil {
					result[url] = e
				} else {
					result[url] = string(b)
				}
			}

		}(url)
	}

	wg.Wait()
	out := n.formatResponse(result)
	return out, nil
}

func (n *NodeService) formatResponse(result map[string]interface{}) string {
	out := ""
	for k, v := range result {

		if k == NODE_MIRROR {
			k = "NodeJS"
		} else {
			k = "IO.JS"
		}

		ve, ok := v.(error)
		if ok {
			fmt.Sprintf("Got error: %v\n", ve)
			continue
		}

		vv, o := v.(string)

		if !o {
			fmt.Printf("No response\n")
			continue
		}

		r, _ := regexp.Compile("<a[\\s\\w=\".\\/]*>([v\\d.]*)\\/\\s*<\\/a>")

		version := r.FindAllStringSubmatch(vv, -1)

		var versions []string

		for _, s := range version {
			if len(s) == 2 && s[1] != ".." {
				versions = append(versions, s[1])
			}
		}

		out = out + fmt.Sprintf("\n%s:\n %s\n", k, versions)
	}

	return out
}

func (n *NodeService) RemoteFile(v string, arch string, platform string) (string, string) {
	var filename string
	var mirror string

	version, prefix := n.GetPrefix(v)

	if prefix == "io" {
		filename = fmt.Sprintf("iojs-%s-%s-%s.tar.gz", version, platform, arch)
		mirror = IOJS_MIRROR
	} else {
		filename = fmt.Sprintf("node-%s-%s-%s.tar.gz", version, platform, arch)
		mirror = NODE_MIRROR
	}

	return mirror + version + "/" + filename, filename
}

func (n *NodeService) GetPrefix(version string) (v string, p string) {
	r, _ := regexp.Compile("^([a-zA-Z]*)@.*")

	prefix := r.FindStringSubmatch(version)

	if len(prefix) == 2 {
		p = prefix[1]
		v = strings.Replace(version, prefix[1]+"@", "", 1)
	} else {
		p = "node"
		v = version
	}

	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}

	return
}
