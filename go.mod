module github.com/dafiti-group/k8s-values-updater

go 1.14

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-git/go-git/v5 v5.1.0
	github.com/go-logr/logr v0.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/kustomize/kyaml v0.3.4
)
