package models

type ServerAppSpecModel struct {
	Version float32 `json:"version" yaml:"version"`
	OS      string  `json:"os" yaml:"os"`
	Files   []File  `json:"files" yaml:"files"`

	// Optional
	Permissions []Permission      `json:"permissions" yaml:"permissions"`
	Hooks       map[string][]Hook `json:"hooks" yaml:"hooks"`
}

type File struct {
	Source      string `json:"source" yaml:"source"`
	Destination string `json:"destination" yaml:"destination"`
}

type Permission struct {
	Object  string   `json:"object" yaml:"object"`
	Pattern string   `json:"pattern" yaml:"pattern"`
	Except  string   `json:"except" yaml:"except"`
	Owner   string   `json:"owner" yaml:"owner"`
	Group   string   `json:"group" yaml:"group"`
	Mode    string   `json:"mode" yaml:"mode"`
	Acls    []string `json:"acls" yaml:"acls"`
	Context Context  `json:"context" yaml:"context"`
	Type    []string `json:"type" yaml:"type"`
}

type Context struct {
	User  string `json:"user" yaml:"user"`
	Type  string `json:"type" yaml:"type"`
	Range string `json:"range" yaml:"range"`
}

type Hook struct {
	Location string `json:"location" yaml:"location"`
	Timeout  string `json:"timeout" yaml:"timeout"`
	Runas    string `json:"runas" yaml:"runas"`
}
