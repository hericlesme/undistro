/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_generateYAML(t *testing.T) {
	g := NewWithT(t)
	// create a local template
	contents := `v1: ${VAR1:=default1}
v2: ${VAR2=default2}
v3: ${VAR3:-default3}`
	template, cleanup1 := createTempFile(g, contents)
	defer cleanup1()

	templateWithoutVars, cleanup2 := createTempFile(g, `v1: foobar
v2: bazfoo`)
	defer cleanup2()

	inputReader := strings.NewReader(contents)

	tests := []struct {
		name           string
		options        *generateYAMLOptions
		inputReader    io.Reader
		expectErr      bool
		expectedOutput string
	}{
		{
			name:      "prints processed yaml using --from flag",
			options:   &generateYAMLOptions{url: template},
			expectErr: false,
			expectedOutput: `v1: default1
v2: default2
v3: default3
`,
		},
		{
			name:      "prints variables using --list-variables flag",
			options:   &generateYAMLOptions{url: template, listVariables: true},
			expectErr: false,
			expectedOutput: `Variables:
  - VAR1
  - VAR2
  - VAR3
`,
		},
		{
			name:      "returns error for bad templateFile path",
			options:   &generateYAMLOptions{url: "/tmp/do-not-exist", listVariables: true},
			expectErr: true,
		},
		{
			name:      "returns error if no options were specified",
			options:   &generateYAMLOptions{},
			expectErr: true,
		},
		{
			name:           "prints nothing if there are no variables in the template",
			options:        &generateYAMLOptions{url: templateWithoutVars, listVariables: true},
			expectErr:      false,
			expectedOutput: "\n",
		},
		{
			name:        "prints processed yaml using specified reader when '--from=-'",
			options:     &generateYAMLOptions{url: "-", listVariables: false},
			inputReader: inputReader,
			expectErr:   false,
			expectedOutput: `v1: default1
v2: default2
v3: default3
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			gyOpts = tt.options
			buf := bytes.NewBufferString("")
			err := generateYAML(inputReader, buf)
			if tt.expectErr {
				g.Expect(err).To(HaveOccurred())
				return
			}

			output, err := ioutil.ReadAll(buf)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(string(output)).To(Equal(tt.expectedOutput))
		})
	}

}

// createTempFile creates a temporary yaml file inside a temp dir. It returns
// the filepath and a cleanup function for the temp directory.
func createTempFile(g *WithT, contents string) (string, func()) {
	dir, err := ioutil.TempDir("", "undistro")
	g.Expect(err).NotTo(HaveOccurred())

	templateFile := filepath.Join(dir, "templ.yaml")
	g.Expect(ioutil.WriteFile(templateFile, []byte(contents), 0600)).To(Succeed())

	return templateFile, func() {
		os.RemoveAll(dir)
	}
}
