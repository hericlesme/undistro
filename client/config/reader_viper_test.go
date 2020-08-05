/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_viperReader_Init(t *testing.T) {
	g := NewWithT(t)

	// Change HOME dir and do not specify config file
	// (.undistro/undistro) in it.
	undistroHomeDir, err := ioutil.TempDir("", "undistro-default")
	g.Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(undistroHomeDir)

	dir, err := ioutil.TempDir("", "undistro")
	g.Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(dir)

	configFile := filepath.Join(dir, "undistro.yaml")
	g.Expect(ioutil.WriteFile(configFile, []byte("bar: bar"), 0600)).To(Succeed())

	configFileBadContents := filepath.Join(dir, "undistro-bad.yaml")
	g.Expect(ioutil.WriteFile(configFileBadContents, []byte("bad-contents"), 0600)).To(Succeed())

	tests := []struct {
		name       string
		configPath string
		configDirs []string
		expectErr  bool
	}{
		{
			name:       "reads in config successfully",
			configPath: configFile,
			configDirs: []string{undistroHomeDir},
			expectErr:  false,
		},
		{
			name:       "returns error for invalid config file path",
			configPath: "do-not-exist.yaml",
			configDirs: []string{undistroHomeDir},
			expectErr:  true,
		},
		{
			name:       "does not return error if default file doesn't exist",
			configPath: "",
			configDirs: []string{undistroHomeDir},
			expectErr:  false,
		},
		{
			name:       "returns error for malformed config",
			configPath: configFileBadContents,
			configDirs: []string{undistroHomeDir},
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gg := NewWithT(t)
			v := newViperReader(InjectConfigPaths(tt.configDirs))
			if tt.expectErr {
				gg.Expect(v.Init(tt.configPath)).ToNot(Succeed())
				return
			}
			gg.Expect(v.Init(tt.configPath)).To(Succeed())

		})
	}
}

func Test_viperReader_Get(t *testing.T) {
	g := NewWithT(t)

	dir, err := ioutil.TempDir("", "undistro")
	g.Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(dir)

	os.Setenv("FOO", "foo")

	configFile := filepath.Join(dir, "undistro.yaml")
	g.Expect(ioutil.WriteFile(configFile, []byte("bar: bar"), 0600)).To(Succeed())

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Read from env",
			args: args{
				key: "FOO",
			},
			want:    "foo",
			wantErr: false,
		},
		{
			name: "Read from file",
			args: args{
				key: "BAR",
			},
			want:    "bar",
			wantErr: false,
		},
		{
			name: "Fails if missing",
			args: args{
				key: "BAZ",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := NewWithT(t)

			v := newViperReader(InjectConfigPaths([]string{dir}))

			gs.Expect(v.Init(configFile)).To(Succeed())

			got, err := v.Get(tt.args.key)
			if tt.wantErr {
				gs.Expect(err).To(HaveOccurred())
				return
			}

			gs.Expect(err).NotTo(HaveOccurred())
			gs.Expect(got).To(Equal(tt.want))
		})
	}
}

func Test_viperReader_GetWithoutDefaultConfig(t *testing.T) {
	g := NewWithT(t)
	dir, err := ioutil.TempDir("", "undistro")
	g.Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(dir)

	os.Setenv("FOO_FOO", "bar")

	v := newViperReader(InjectConfigPaths([]string{dir}))
	g.Expect(v.Init("")).To(Succeed())

	got, err := v.Get("FOO_FOO")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(got).To(Equal("bar"))
}

func Test_viperReader_Set(t *testing.T) {
	g := NewWithT(t)

	dir, err := ioutil.TempDir("", "undistro")
	g.Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(dir)

	os.Setenv("FOO", "foo")

	configFile := filepath.Join(dir, "undistro.yaml")

	g.Expect(ioutil.WriteFile(configFile, []byte("bar: bar"), 0600)).To(Succeed())

	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				key:   "FOO",
				value: "bar",
			},
			want: "bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := NewWithT(t)

			v := &viperReader{}

			gs.Expect(v.Init(configFile)).To(Succeed())

			v.Set(tt.args.key, tt.args.value)

			got, err := v.Get(tt.args.key)
			gs.Expect(err).NotTo(HaveOccurred())
			gs.Expect(got).To(Equal(tt.want))
		})
	}
}

func Test_viperReader_checkDefaultConfig(t *testing.T) {
	g := NewWithT(t)
	dir, err := ioutil.TempDir("", "undistro")
	g.Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(dir)
	dir = strings.TrimSuffix(dir, "/")

	configFile := filepath.Join(dir, "undistro.yaml")
	g.Expect(ioutil.WriteFile(configFile, []byte("bar: bar"), 0600)).To(Succeed())

	type fields struct {
		configPaths []string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "tmp path without final /",
			fields: fields{
				configPaths: []string{dir},
			},
			want: true,
		},
		{
			name: "tmp path with final /",
			fields: fields{
				configPaths: []string{fmt.Sprintf("%s/", dir)},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := NewWithT(t)

			v := &viperReader{
				configPaths: tt.fields.configPaths,
			}
			gs.Expect(v.checkDefaultConfig()).To(Equal(tt.want))
		})
	}
}
