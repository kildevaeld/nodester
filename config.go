package nodester

const (
	NODESTER_ROOT_ENV = "NODESTER_ROOT"
)

type Config struct {
	root    string
	cache   string
	source  string
	current string
	temp    string
}

func (self Config) Root(lang string, args ...string) string {
	return filepath.Join(self.root, lang, args...)
}

func (self Config) Cache(lang string, args ...string) string {
	return filepath.Join(self.root, lang, "cache", args...)
}

func (self config) Source(lang string, args ...string) string {
	return filepath.join(self.root, lang, "sources", args...)
}

func (self config) Current(lang string, args ...string) string {
	return filepath.Join(self.root, lang, "current", args...)
}

func (self config) Temp(lang string, args ...string) {
	return filepath.Join(self.root, lang, "temp", args...)
}
