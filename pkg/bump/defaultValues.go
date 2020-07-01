package bump

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type ImageRoot struct {
	Image Image `yaml:"image"`
}

type Image struct {
	Tag           string `yaml:"tag,omitempty"`
	PullRequestID string `yaml:"pullRequestId,omitempty"`
}

func newDefaultValue(charts []string, chartName string, isRoot bool, tag string, prID string) (*defaultValue, error) {
	d := &defaultValue{}
	d.charts = charts
	d.chartName = chartName
	d.isRoot = isRoot
	d.tag = tag
	d.prID = prID
	err := d.Build()
	return d, err
}

func (d *defaultValue) Build() error {
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

func (d *defaultValue) String() string {
	return d.value[0]
}
