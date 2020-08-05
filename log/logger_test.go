/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package log

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/pkg/errors"
)

func TestFlatten(t *testing.T) {
	type args struct {
		prefix string
		kvList []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "message without values",
			args: args{
				prefix: "",
				kvList: []interface{}{
					"msg", "this is a message",
				},
			},
			want: "this is a message",
		},
		{
			name: "message with values",
			args: args{
				prefix: "",
				kvList: []interface{}{
					"msg", "this is a message",
					"val1", 123,
					"val2", "string",
					"val3", "string with spaces",
				},
			},
			want: "this is a message val1=123 val2=\"string\" val3=\"string with spaces\"",
		},
		{
			name: "error without values",
			args: args{
				prefix: "",
				kvList: []interface{}{
					"msg", "this is a message",
					"error", errors.New("this is an error"),
				},
			},
			want: "this is a message: this is an error",
		},
		{
			name: "error with values",
			args: args{
				prefix: "",
				kvList: []interface{}{
					"msg", "this is a message",
					"error", errors.New("this is an error"),
					"val1", 123,
				},
			},
			want: "this is a message: this is an error val1=123",
		},
		{
			name: "message with prefix",
			args: args{
				prefix: "a\\b",
				kvList: []interface{}{
					"msg", "this is a message",
				},
			},
			want: "[a\\b] this is a message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got, err := flatten(logEntry{
				Prefix: tt.args.prefix,
				Level:  0,
				Values: tt.args.kvList,
			})
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(got).To(Equal(tt.want))
		})
	}
}
