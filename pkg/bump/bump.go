package bump

import (
	"fmt"

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
	token string,
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

	git.Branch = "feature/auth-only-with-https"
	git.URL = "https://github.com/dafiti-group/k8s-values-updater.git"

	separator := ","
	err := git.Init(token, b.FileNames, b.DirPath, separator)
	if err != nil {
		return err
	}

	err = file.Init()
	if err != nil {
		return err
	}

	return nil
}

// Run ...
func (b *Bump) Run() error {
	//
	if err := b.git.Sync(); err != nil {
		return err
	}

	//
	files := b.git.Files()

	if len(files) < 1 {
		return fmt.Errorf("file not found")
	}

	//
	for _, f := range files {
		if err := b.file.Bump(f); err != nil {
			return err
		}
	}

	if !b.file.HasChanges() {
		b.Log.Info("nothing changed")
		return nil
	}

	if err := b.git.Push(b.file.GetChanges()); err != nil {
		return err
	}

	return nil
}
