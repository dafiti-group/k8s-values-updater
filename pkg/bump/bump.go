package bump

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/k0kubun/pp"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
	kyamlmerge "sigs.k8s.io/kustomize/kyaml/yaml/merge2"
)

var directory = "."

func (b *Bump) Init(user string, pass string, dryRun bool) error {
	// TODO: Validade fields
	b.DryRun = dryRun

	err := b.SetBasicAuth(user, pass)
	if err != nil {
		return err
	}
	return nil
}

// https://github.com/divramod/dp/blob/master/utils/git/main.go
func (b *Bump) Files() ([]string, error) {
	r := []string{}

	for _, v := range strings.Split(b.FileNames, ",") {
		f := filepath.Join(b.DirPath, v)
		matches, err := filepath.Glob(f)
		if err != nil {
			return nil, err
		}
		for _, v = range matches {
			r = append(r, v)
		}
	}
	fmt.Println(r)
	return r, nil
}

func (b *Bump) Run() error {
	files, err := b.Files()
	if err != nil {
		return err
	}
	//
	pp.Println(files)
	if len(files) < 1 {
		return fmt.Errorf("File not found")
	}
	//
	if err = b.sync(); err != nil {
		return err
	}
	//
	for _, v := range files {
		if err = b.bump(v); err != nil {
			return err
		}
	}

	//
	if b.HasNoChanges() {
		fmt.Println("Nothing Changed")
		return nil
	}

	//
	if err = b.push(files); err != nil {
		return err
	}

	return nil
}
func (b *Bump) sync() error {
	// Opens an already existing repository.
	r, err := git.PlainOpen(directory)
	if err != nil {
		return err
	}

	// Get WorkTree
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// Pull
	fmt.Println("Pull")
	err = w.Pull(&git.PullOptions{
		RemoteName: b.RemoteName,
		Force:      true,
		Auth:       b.auth,
	})
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (b *Bump) push(files []string) error {
	//
	commitMsg := "[ci skip] ci: edit values with the new image tag\n\n\nskip-checks: true"
	name := "K8s Values Updater"

	// Opens an already existing repository.
	r, err := git.PlainOpen(directory)
	if err != nil {
		return err
	}

	// Get WorkTree
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// Add
	for _, f := range files {
		_, err = w.Add(f)
		if err != nil {
			return err
		}
	}

	// Status
	status, err := w.Status()
	if err != nil {
		return err
	}
	fmt.Println(status)

	// Create Commit
	commit, err := w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: b.Email,
			When:  time.Now(),
		},
	})

	// Commit
	obj, err := r.CommitObject(commit)
	if err != nil {
		return err
	}
	fmt.Println(obj)

	// Push
	fmt.Println("Push")
	err = r.Push(&git.PushOptions{
		RemoteName: b.RemoteName,
		Auth:       b.auth,
	})
	if err != nil {
		return err
	}

	//
	fmt.Println("Ok")
	return nil
}

func (b *Bump) bump(filePath string) error {
	o, err := kyaml.ReadFile(filePath)
	if err != nil {
		return err
	}

	charts, err := o.Fields()
	if err != nil {
		return err
	}

	values, err := newDefaultValue(charts, b.ChartName, b.IsRoot, b.Tag, b.PrID)
	if err != nil {
		return err
	}

	buffer := values.String()
	if b.ReplaceWith != "" {
		buffer = b.ReplaceWith
	}
	src, err := kyaml.Parse(buffer)
	if err != nil {
		return err
	}

	b.Before, err = o.String()
	if err != nil {
		return nil
	}

	_, err = kyamlmerge.Merge(src, o)
	if err != nil {
		return err
	}

	b.After, err = o.String()
	if err != nil {
		return nil
	}

	if b.HasNoChanges() {
		fmt.Println("Nothing Changed")
		return nil
	}

	err = kyaml.WriteFile(o, filePath)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bump) HasNoChanges() bool {
	return b.Before == b.After
}

func (b *Bump) SetBasicAuth(user string, pass string) error {
	if pass != "" && !b.ForceSSH {
		fmt.Println("Will authenticate with basic auth")
		b.auth = &http.BasicAuth{
			Username: user,
			Password: pass,
		}
	} else {
		fmt.Println("Will authenticate with SSH")
		sshPath := path.Join(os.Getenv("HOME"), ".ssh", "id_rsa")

		sshKey, err := ioutil.ReadFile(sshPath)
		if err != nil {
			return err
		}

		key, err := ssh.NewPublicKeys("git", []byte(sshKey), "")
		if err != nil {
			return err
		}
		b.auth = key
	}

	return nil
}
