package models

type Ec2OnPremAppSpecModel struct {
	Version     string
	OS          string
	Files       []File
	Permissions []Permission
	Hooks       map[string]Hook
}

type File struct {
	Source      string
	Destination string
}

type Permission struct {
	Object     string
	Pattern    string
	Except     string
	Owner      string
	Group      string
	Mode       string
	Acls       []string
	ContextObj []Context
	Type       []string
}

type Context struct {
	User  string
	Type  string
	Range string
}

type Hook struct {
	Location string
	Timeout  string
	Runas    string
}
