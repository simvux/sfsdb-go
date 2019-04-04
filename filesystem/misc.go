package filesystem

import (
	"strings"
)

// Filepath is a safe way to construct file paths that cannot be escaped
type Filepath []string

func NewFilepath(base string) Filepath {
	return Filepath{base}
}

func (p *Filepath) Append(s string) {
	sanatize(&s)
	*p = append(*p, s)
}

func (p *Filepath) Unwrap() string {
	return strings.Join(*p, "/")
}

func sanatize(s *string) {
	r := strings.NewReplacer("/", `\/`, `\`, `\\`)
	*s = r.Replace(*s)
}
