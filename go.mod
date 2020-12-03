module github.com/getupio-undistro/undistro

go 1.15

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/go-logr/logr v0.3.0
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/pkg/errors v0.9.1
	helm.sh/helm/v3 v3.4.1
	k8s.io/api v0.19.4
	k8s.io/apiextensions-apiserver v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/cli-runtime v0.19.4
	k8s.io/client-go v0.19.4
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.7.0-alpha.8
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)
