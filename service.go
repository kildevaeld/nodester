package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type Service interface {
	Name() string
	ListRemote() (string, error)
	Latest(string, string, string) (string, error)
	RemoteFile(version string, arch string, platform string) (string, string)
	GetPrefix(str string) (version string, prefix string)
}

type NodeService struct{}

func (n *NodeService) Name() string { return "node" }

func (n *NodeService) ListRemote() (string, error) {
	var lock sync.Mutex
	result := make(map[string]interface{})

	var wg sync.WaitGroup

	urls := []string{IOJS_MIRROR}

	for _, url := range urls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			res, err := http.Get(url)

			if err != nil {
				fmt.Printf("Error %s\n", err.Error())
				return
			}

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

func (n *NodeService) Latest(prefix string, platform string, arch string) (string, error) {
	var mirror string
	if prefix == "io" {
		mirror = IOJS_MIRROR
	} else {
		mirror = NODE_MIRROR
	}
	mirror = mirror + "latest"

	resp, err := http.Get(mirror)

	if err != nil {
		return "", nil
	}

	body, _ := ioutil.ReadAll(resp.Body)

	r, _ := regexp.Compile("<a[\\s\\w=\".\\/-]*>[v\\w]+-([-\\d.v]*)")

	versions := r.FindAllStringSubmatch(string(body), -1)

	var version string
	for _, v := range versions {
		if strings.Contains(v[0], platform+"-"+arch) {
			version = strings.Trim(v[1], "-")
		}
	}

	return version, nil

}

func (n *NodeService) formatResponse(result map[string]interface{}) string {
	out := ""
	for k, v := range result {

		if k == NODE_MIRROR {
			k = "node.js"
		} else {
			k = "io.js"
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

		var versions sort.StringSlice

		for _, s := range version {
			if len(s) == 2 && s[1] != ".." {
				if s[1][0] == 'v' {
					s[1] = s[1][1:]
				}
				versions = append(versions, s[1])
			}
		}
		sort.Strings(versions)
		sort.Sort(sort.Reverse(versions))

		fVersions := strings.Join(versions[0:10], "\t")
		out = out + fmt.Sprintf("\n  %s:\n   %s\n", k, fVersions)
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

	if !strings.HasPrefix(v, "v") && v != "latest" {
		v = "v" + v
	}

	return
}
