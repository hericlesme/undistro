/*
Copyright 2020 The UnDistro authors

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
package cli

import (
	"flag"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ConfigFlags struct {
	ConfigFile *string
	Verbosity  *int
	KubeFlags  *genericclioptions.ConfigFlags
}

// NewConfigFlags returns ConfigFlags with default values set
func NewConfigFlags() *ConfigFlags {
	home := homedir.HomeDir()
	cfgPath := filepath.Join(home, ".undistro", "undistro.yaml")
	return &ConfigFlags{
		Verbosity:  intptr(0),
		ConfigFile: stringptr(cfgPath),
		KubeFlags:  genericclioptions.NewConfigFlags(true),
	}
}

func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet, goflags *flag.FlagSet) {
	if f.Verbosity != nil {
		goflags.IntVar(f.Verbosity, "v", *f.Verbosity, "Set the log level verbosity.")
	}
	flags.AddGoFlagSet(goflags)
	if f.ConfigFile != nil {
		flags.StringVar(f.ConfigFile, "config", *f.ConfigFile, "Path to undistro configuration (default is `$HOME/.undistro/undistro.yaml`)")
	}
	f.KubeFlags.AddFlags(flags)
}

func (f *ConfigFlags) Init() func() {
	return func() {
		log.SetLogger(log.Log.V(*f.Verbosity))
		viper.SetConfigFile(*f.ConfigFile)
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err == nil {
			log.Log.Info("using config", "file", viper.ConfigFileUsed())
		}
	}
}

func stringptr(s string) *string {
	return &s
}

func intptr(i int) *int {
	return &i
}
