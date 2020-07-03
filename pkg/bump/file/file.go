package file

import (
	"github.com/sirupsen/logrus"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
	kyamlmerge "sigs.k8s.io/kustomize/kyaml/yaml/merge2"
)

type File struct {
	IsRoot      bool
	ChartName   string
	ReplaceWith string
	Tag         string
	PrID        string
	before      string
	after       string
	Log         *logrus.Logger
}

//
func (f *File) Bump(filePath string) error {
	// Read file
	o, err := kyaml.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Get Chart root fields
	charts, err := o.Fields()
	if err != nil {
		return err
	}

	// Build the values to do the replace
	values, err := CreateDefaultValuesString(charts, f.ChartName, f.IsRoot, f.Tag, f.PrID, f.ReplaceWith)
	if err != nil {
		return err
	}

	// Parse into a kuztomize yaml
	src, err := kyaml.Parse(values)
	if err != nil {
		return err
	}

	// BKP For future compare
	f.before, err = o.String()
	if err != nil {
		return nil
	}

	// Merge
	_, err = kyamlmerge.Merge(src, o)
	if err != nil {
		return err
	}

	// Save resulr for future compare
	f.after, err = o.String()
	if err != nil {
		return nil
	}

	// If Nothing Changed do nothinb
	if f.HasNoChanges() {
		f.Log.Info("nothing changed, no file will be altered")
		return nil
	}

	// Save File
	err = kyaml.WriteFile(o, filePath)
	if err != nil {
		return err
	}

	return nil
}

// HasNoChanges Return true if nothing in the file changed
func (f *File) HasNoChanges() bool {
	return f.before == f.after
}
