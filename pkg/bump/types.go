package bump

import "github.com/go-git/go-git/v5/plumbing/transport"

type Bump struct {
	DirPath     string
	FileNames   string
	RemoteName  string
	RefSpec     string
	Branch      string
	Email       string
	IsRoot      bool
	ChartName   string
	DryRun      bool
	ReplaceWith string
	Tag         string
	PrID        string
	Before      string
	After       string
	auth        transport.AuthMethod
	ForceSSH    bool
}

type defaultValue struct {
	charts    []string
	chartName string
	isRoot    bool
	tag       string
	prID      string
	value     []string
}
