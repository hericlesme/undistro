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
package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

var (
	cfgFile   string
	verbosity *int
)

var RootCmd = &cobra.Command{
	Use:          "undistro",
	SilenceUsage: true,
	Short:        "undistro controls the unDistro kubernetes distribution",
	Long:         LongDesc(`undistro controls the unDistro kubernetes distribution`),
}

func Execute(ctx context.Context) {
	if err := RootCmd.ExecuteContext(ctx); err != nil {
		if verbosity != nil && *verbosity >= 5 {
			if err, ok := err.(stackTracer); ok {
				for _, f := range err.StackTrace() {
					fmt.Fprintf(os.Stderr, "%+s:%d\n", f, f)
				}
			}
		}
		os.Exit(1)
	}
}

func init() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	verbosity = flag.CommandLine.Int("v", 0, "Set the log level verbosity.")
	flags := RootCmd.PersistentFlags()
	flags.AddGoFlagSet(flag.CommandLine)
	flags.StringVar(&cfgFile, "config", "", "Path to undistro configuration (default is `$HOME/.undistro/undistro.yaml`)")
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	kubeConfigFlags.AddFlags(flags)
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	log.SetLogger(log.Log.V(*verbosity))
}

const Indentation = `  `

// LongDesc normalizes a command's long description to follow the conventions.
func LongDesc(s string) string {
	if len(s) == 0 {
		return s
	}
	return normalizer{s}.heredoc().trim().string
}

// Examples normalizes a command's examples to follow the conventions.
func Examples(s string) string {
	if len(s) == 0 {
		return s
	}
	return normalizer{s}.trim().indent().string
}

// TODO: document this, what does it do? Why is it here?
type normalizer struct {
	string
}

func (s normalizer) heredoc() normalizer {
	s.string = heredoc.Doc(s.string)
	return s
}

func (s normalizer) trim() normalizer {
	s.string = strings.TrimSpace(s.string)
	return s
}

func (s normalizer) indent() normalizer {
	splitLines := strings.Split(s.string, "\n")
	indentedLines := make([]string, 0, len(splitLines))
	for _, line := range splitLines {
		trimmed := strings.TrimSpace(line)
		indented := Indentation + trimmed
		indentedLines = append(indentedLines, indented)
	}
	s.string = strings.Join(indentedLines, "\n")
	return s
}
