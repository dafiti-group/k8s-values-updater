package file

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type defaultValue struct {
	charts    []string
	chartName string
	isRoot    bool
	tag       string
	prID      string
	value     []string
}

type ImageRoot struct {
	Image Image `yaml:"image"`
}

type Image struct {
	Tag           string `yaml:"tag,omitempty"`
	PullRequestID string `yaml:"pullRequestId,omitempty"`
}

func CreateDefaultValuesString(charts []string, chartName string, isRoot bool, tag string, prID string, replaceWith string) (string, error) {
	if replaceWith != "" {
		return replaceWith, nil
	}
	d := &defaultValue{}
	d.charts = charts
	d.chartName = chartName
	d.isRoot = isRoot
	d.tag = tag
	d.prID = prID
	err := d.build()
	if err != nil {
		return "", err
	}

	return d.string(), err
}

func (d *defaultValue) build() error {
	image := ImageRoot{
		Image: Image{
			Tag:           d.tag,
			PullRequestID: d.prID,
		},
	}
	//
	if d.chartName != "" {
		buffer, err := yaml.Marshal(map[string]ImageRoot{
			d.chartName: image,
		})
		if err != nil {
			return err
		}
		d.value = []string{string(buffer)}
		return nil
	}
	//
	if d.isRoot {
		buffer, err := yaml.Marshal(image)
		if err != nil {
			return err
		}
		d.value = []string{string(buffer)}
		return nil
	}

	//
	for _, chartName := range d.charts {
		buffer, err := yaml.Marshal(map[string]ImageRoot{
			chartName: image,
		})
		if err != nil {
			return err
		}
		d.value = append(d.value, string(buffer))
	}

	//
	return nil
}

func (d *defaultValue) string() string {
	return d.value[0]
}
