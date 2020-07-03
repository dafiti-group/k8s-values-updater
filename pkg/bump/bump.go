package bump

import (
	"fmt"
	"path/filepath"
	"strings"

	file "github.com/dafiti-group/k8s-values-updater/pkg/bump/file"
	git "github.com/dafiti-group/k8s-values-updater/pkg/bump/git"
)

type Bump struct {
	DirPath   string
	FileNames string
	DryRun    bool

	git  git.Git
	file file.File
}

// Run
func (b *Bump) Init(g *git.Git, file *file.File, user string, pass string, dryRun bool) error {
	// TODO: Validade fields
	b.DryRun = dryRun

	err := g.SetBasicAuth(user, pass)
	if err != nil {
		return err
	}
	return nil
}

// Run
func (b *Bump) Run() error {
	files, err := b.files()
	if err != nil {
		return err
	}
	//
	if len(files) < 1 {
		return fmt.Errorf("File not found")
	}

	//
	if err = b.git.Sync(); err != nil {
		return err
	}

	//
	for _, f := range files {
		if err = b.file.Bump(f); err != nil {
			return err
		}
	}

	//
	if b.file.HasNoChanges() {
		fmt.Println("Nothing Changed")
		return nil
	}

	//
	if err = b.git.Push(files); err != nil {
		return err
	}

	return nil
}

// https://github.com/divramod/dp/blob/master/utils/git/main.go
func (b *Bump) files() ([]string, error) {
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
