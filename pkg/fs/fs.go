/*
Copyright 2020-2021 The UnDistro authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package fs

import (
	"embed"
)

//go:embed clustertemplates
var FS embed.FS

//go:generate helm package -u ../../charts/metallb -d chart
//go:generate helm package -u ../../charts/cert-manager -d chart
//go:generate helm package -u ../../charts/cluster-api -d chart
//go:generate helm package -u ../../charts/ingress-nginx -d chart
//go:generate helm package -u ../../charts/undistro -d chart
//go:generate helm package -u ../../charts/undistro-aws -d chart
//go:generate helm package -u ../../charts/undistro-openstack -d chart
//go:embed chart
var ChartFS embed.FS

//go:embed defaultarch
var DefaultArchFS embed.FS

//go:embed policies
var PoliciesFS embed.FS
