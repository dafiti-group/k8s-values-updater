package bump

import (
	"fmt"
	"path/filepath"
	"strings"

	file "github.com/dafiti-group/k8s-values-updater/pkg/bump/file"
	git "github.com/dafiti-group/k8s-values-updater/pkg/bump/git"
	"github.com/sirupsen/logrus"
)

type Bump struct {
	DirPath   string
	FileNames string
	DryRun    bool

	git  *git.Git
	file *file.File
	Log  *logrus.Logger
}

// Init ...
func (b *Bump) Init(
	git *git.Git,
	file *file.File,
	user string,
	pass string,
	dryRun bool,
	log *logrus.Logger,
) error {

	// Initizali logger
	git.Log = log
	file.Log = log
	b.Log = log

	b.git = git
	b.file = file
	// TODO: Validade fields
	b.DryRun = dryRun

	err := git.SetBasicAuth(user, pass)
	if err != nil {
		return err
	}
	return nil
}

// Run ...
func (b *Bump) Run() error {
	files, err := b.files()
	if err != nil {
		return err
	}
	//
	if len(files) < 1 {
		return fmt.Errorf("file not found")
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
		b.Log.Info("nothing Changed will not push")
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
	b.Log.Debug(r)
	return r, nil
}
