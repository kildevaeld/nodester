package nodester

type Definition struct {
    Name string
    Description string
}

type paths struct {
	root    string
}

func (self paths) root(args ...string) string {
	return filepath.Join(self.root,args...)
}

func (self paths) Cache(args ...string) string {
	return filepath.Join(self.root, "cache", args...)
}

func (self paths) Source(args ...string) string {
	return filepath.join(self.root, "sources", args...)
}

func (self paths) Current(args ...string) string {
	return filepath.Join(self.root, "current", args...)
}

func (self paths) Temp(args ...string) {
	return filepath.Join(self.root, "temp", args...)
}

struct Language {
    paths paths
    definition Definition
}


func NewLanguage(string path, def Definition) (*Language, error) {
   
   lang := &Language{paths{path},def} 
   
   err := ensureDir(path)
   
   return lang, err
}