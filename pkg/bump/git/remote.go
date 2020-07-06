package git

import (
	"fmt"
	"net/url"
)

type Remote struct {
	Repo string   `mapstructure:"repo"`
	Org  string   `mapstructure:"org"`
	url  *url.URL `mapstructure:"url"`
}

// SetUrl ...
func (r *Remote) SetUrl() (err error) {
	if r.Org == "" {
		return fmt.Errorf("org not set")
	}
	if r.Repo == "" {
		return fmt.Errorf("repo not set")
	}

	root := fmt.Sprintf("https://github.com/%v/%v.git", r.Org, r.Repo)
	u, err := url.Parse(root)
	if err != nil {
		return err
	}
	u.Scheme = "https"
	u.Host = "github.com"
	r.url = u
	return err
}

// GetUrl ...
func (r *Remote) GetUrl() string {
	return r.url.String()
}
