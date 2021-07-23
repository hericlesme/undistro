module github.com/getupio-undistro/undistro

go 1.16

replace sigs.k8s.io/cluster-api => github.com/getupio-undistro/cluster-api v0.3.11-0.20210608202222-4052f5f87b90

require (
	github.com/MakeNowJust/heredoc/v2 v2.0.1
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/aws/aws-sdk-go v1.40.6
	github.com/go-logr/logr v0.4.0
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0
	github.com/gorilla/mux v1.8.0
	github.com/json-iterator/go v1.1.11
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0
	github.com/pkg/errors v0.9.1
	github.com/smallstep/truststore v0.9.6
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	helm.sh/helm/v3 v3.6.3
	k8s.io/api v0.21.3
	k8s.io/apiextensions-apiserver v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/cli-runtime v0.21.3
	k8s.io/client-go v0.21.3
	k8s.io/klog/v2 v2.10.0
	k8s.io/kubectl v0.21.3
	k8s.io/utils v0.0.0-20210722164352-7f3ee0f31471
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/cluster-api v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.9.3
	sigs.k8s.io/yaml v1.2.0
)
