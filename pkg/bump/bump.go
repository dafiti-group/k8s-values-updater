package bump

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/k0kubun/pp"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
	kyamlmerge "sigs.k8s.io/kustomize/kyaml/yaml/merge2"
)

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
	if len(files) < 1 {
		return fmt.Errorf("File not found")
	}
	//
	for _, v := range files {
		if err = b.bump(v); err != nil {
			return err
		}
	}

	if b.HasNoChanges() {
		fmt.Println("Nothing Changed")
		return nil
	}

	if err = b.push(files); err != nil {
		return err
	}

	return nil
}
func (b *Bump) push(files []string) error {
	//
	directory := "."
	commitMsg := "[ci skip] ci: edit values with the new image tag\n\n\nskip-checks: true"
	name := "K8s Values Updater"
	key := SSHKeyGet()

	// refSpec := []config.RefSpec{}
	refSpec := []config.RefSpec{
		config.RefSpec("+refs/heads/master:refs/heads/master"),
	}
	// if b.Branch != "" {
	// 	refSpec = []config.RefSpec{config.RefSpec(fmt.Sprintf("+refs/heads/%v:refs/remotes/origin/%v", b.Branch, b.Branch))}
	// }
	// if b.RefSpec != "" {
	// 	refSpec = []config.RefSpec{config.RefSpec(b.RefSpec)}
	// }

	// Opens an already existing repository.
	r, err := git.PlainOpen(directory)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	c, err := r.Config()
	if err != nil {
		return err
	}
	pp.Println(c.Branches)

	// Add
	for _, f := range files {
		_, err = w.Add(f)
		if err != nil {
			return err
		}
	}

	status, err := w.Status()
	if err != nil {
		return err
	}

	fmt.Println(status)
	commit, err := w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: b.Email,
			When:  time.Now(),
		},
	})

	obj, err := r.CommitObject(commit)
	if err != nil {
		return err
	}

	fmt.Println(obj)

	fmt.Println("Fetch")
	err = r.Fetch(&git.FetchOptions{
		RemoteName: b.RemoteName,
		RefSpecs:   refSpec,
	})
	if err != nil {
		return err
	}

	// fmt.Println("Pull")
	// err = w.Pull(&git.PullOptions{
	// 	RemoteName: b.RemoteName,
	// 	Force:      true,
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// }

	fmt.Println("Push")
	err = r.Push(&git.PushOptions{
		RemoteName: b.RemoteName,
		RefSpecs:   refSpec,
		Auth:       key,
	})
	if err != nil {
		return err
	}

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

func SSHKeyGet() *ssh.PublicKeys {
	var publicKey *ssh.PublicKeys
	sshPath := os.Getenv("HOME") + "/.ssh/id_rsa"
	sshKey, _ := ioutil.ReadFile(sshPath)
	publicKey, keyError := ssh.NewPublicKeys("git", []byte(sshKey), "")
	if keyError != nil {
		fmt.Println(keyError)
	}
	return publicKey
}
