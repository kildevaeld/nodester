package nodester

type Manifests []Manifest

type Manifest struct {
	Version string
	Date    string
	Files   []string
	Lts     interface{}
	/*Modules string
	Npm     string
	V8      string
	Uv      string
	Zlib    string
	OpenSSL string
	Lts     bool*/
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
