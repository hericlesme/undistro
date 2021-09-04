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
package cli

import (
	"os"
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_getIPFromConfig(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get ingress address",
			args: args{
				file: "./testdata/config.yaml",
			},
			want: "k8s.felipeweb.dev",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			byt, err := os.ReadFile(tt.args.file)
			if err != nil {
				t.Error(err)
				return
			}
			cfg := make(map[string]interface{})
			err = yaml.Unmarshal(byt, &cfg)
			if err != nil {
				t.Error(err)
				return
			}
			if got := getIPFromConfig(cfg); got != tt.want {
				t.Errorf("getIPFromConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
