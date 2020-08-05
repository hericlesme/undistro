/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package client

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_parseProviderName(t *testing.T) {
	type args struct {
		provider string
	}
	tests := []struct {
		name        string
		args        args
		wantName    string
		wantVersion string
		wantErr     bool
	}{
		{
			name: "simple name",
			args: args{
				provider: "provider",
			},
			wantName:    "provider",
			wantVersion: "",
			wantErr:     false,
		},
		{
			name: "name & version",
			args: args{
				provider: "provider:version",
			},
			wantName:    "provider",
			wantVersion: "version",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			gotName, gotVersion, err := parseProviderName(tt.args.provider)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
			g.Expect(gotName).To(Equal(tt.wantName))

			g.Expect(gotVersion).To(Equal(tt.wantVersion))
		})
	}
}
