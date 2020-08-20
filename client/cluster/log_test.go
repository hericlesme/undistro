/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package cluster

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_logStreamer_readLog(t *testing.T) {
	type args struct {
		reader io.ReadCloser
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
	}{
		{
			"read log with success",
			args{
				ioutil.NopCloser(strings.NewReader("log undistro")),
			},
			"log undistro",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			l := &logStreamer{}
			writer := &bytes.Buffer{}
			l.readLog(context.Background(), tt.args.reader, writer)
			g.Expect(writer.String()).To(Equal(tt.wantWriter))
		})
	}
}
