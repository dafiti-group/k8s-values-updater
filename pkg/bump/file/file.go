package file

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
	kyamlmerge "sigs.k8s.io/kustomize/kyaml/yaml/merge2"
)

type File struct {
	IsRoot       bool
	ChartName    string
	ReplaceWith  string
	Tag          string
	PrID         string
	Log          *logrus.Logger
	changedFiles billy.Filesystem
	changes      int
}

//
func (f *File) Init() error {
	f.changedFiles = memfs.New()
	f.changes = 0
	return nil
}

func (f *File) Bump(ioFile billy.File) error {
	// Convert fs file into buffer
	b, err := ioutil.ReadAll(ioFile)
	if err != nil {
		return err
	}

	// Parse File
	parsedFile, err := kyaml.Parse(string(b))
	if err != nil {
		return err
	}

	// Get Chart root fields
	charts, err := parsedFile.Fields()
	if err != nil {
		return err
	}

	// Build the values to do the replace
	values, err := CreateDefaultValuesString(charts, f.ChartName, f.IsRoot, f.Tag, f.PrID, f.ReplaceWith)
	if err != nil {
		return err
	}

	// Parse into a kuztomize yaml
	parsedValues, err := kyaml.Parse(values)
	if err != nil {
		return err
	}

	// BKP For future compare
	before, err := parsedFile.String()
	if err != nil {
		return err
	}

	// Merge
	_, err = kyamlmerge.Merge(parsedValues, parsedFile)
	if err != nil {
		return err
	}

	// Save result for future compare
	after, err := parsedFile.String()
	if err != nil {
		return err
	}

	// If Nothing Changed do nothing
	if before == after {
		f.Log.Info("nothing changed, no file will be altered")
		return nil
	}

	// Create file
	err = util.WriteFile(f.changedFiles, ioFile.Name(), []byte(after), os.ModePerm)
	if err != nil {
		return err
	}

	// Increment a changes helper
	f.changes += 1

	return nil
}

func (f *File) HasChanges() bool {
	return f.changes >= 1
}
func (f *File) GetChanges() billy.Filesystem {
	return f.changedFiles
}
