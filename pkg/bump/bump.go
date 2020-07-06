package bump

import (
	"fmt"

	file "github.com/dafiti-group/k8s-values-updater/pkg/bump/file"
	git "github.com/dafiti-group/k8s-values-updater/pkg/bump/git"
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
)

type Bump struct {
	DirPath   string         `mapstructure:"dir_path"`
	FileNames string         `mapstructure:"file_names"`
	DryRun    bool           `mapstructure:"dry_run"`
	Log       *logrus.Logger `mapstructure:"log"`

	*git.Git
	*file.File
}

func New(log *logrus.Logger) (b *Bump) {
	pp.Println("Hy")
	b = &Bump{
		Git: &git.Git{
			Log: log,
		},
		File: &file.File{
			Log: log,
		},
	}
	return b
}

// Init ...
// TODO: Validade fields
func (b *Bump) Init(
	token string,
	dryRun bool,
	separator string,
) error {
	// Initialize Params from root
	b.DryRun = dryRun

	b.Log.Info(
		"branch", b.Branch,
		"url", b.URL,
	)

	err := b.Git.Init(token, b.FileNames, b.DirPath, separator)
	if err != nil {
		return err
	}

	err = b.File.Init()
	if err != nil {
		return err
	}

	return nil
}

// Run ...
func (b *Bump) Run() error {
	//
	files := b.Git.Files()

	if len(files) < 1 {
		return fmt.Errorf("file not found")
	}

	//
	for _, f := range files {
		if err := b.File.Bump(f); err != nil {
			return err
		}
	}

	if !b.File.HasChanges() {
		b.Log.Info("nothing changed")
		return nil
	}

	if err := b.Git.Push(b.File.GetChanges()); err != nil {
		return err
	}

	return nil
}
