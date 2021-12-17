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
	"flag"
	"io"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/getupio-undistro/undistro/pkg/version"
	pinnipedcmd "github.com/getupio-undistro/undistro/third_party/pinniped/cmd"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apiresources"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/auth"
	"k8s.io/kubectl/pkg/cmd/delete"
	"k8s.io/kubectl/pkg/cmd/describe"
	"k8s.io/kubectl/pkg/cmd/logs"
	"k8s.io/kubectl/pkg/cmd/patch"
	"k8s.io/kubectl/pkg/cmd/rollout"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

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

func NewUndistroCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "undistro",
		SilenceUsage: true,
		Short:        "undistro controls the unDistro kubernetes distribution",
		Long:         LongDesc(`undistro controls the unDistro kubernetes distribution`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags := cmd.PersistentFlags()
	cfgFlags := NewConfigFlags()
	cfgFlags.AddFlags(flags, flag.CommandLine)
	ioStreams := genericclioptions.IOStreams{In: in, Out: out, ErrOut: err}
	cmd.AddCommand(NewCmdSetup(cfgFlags, ioStreams))
	cmd.AddCommand(NewCmdDestroy(ioStreams))
	f := cmdutil.NewFactory(cfgFlags)
	cmd.AddCommand(auth.NewCmdAuth(f, ioStreams))
	cmd.AddCommand(delete.NewCmdDelete(f, ioStreams))
	cmd.AddCommand(patch.NewCmdPatch(f, ioStreams))
	cmd.AddCommand(apply.NewCmdApply("undistro", f, ioStreams))
	cmd.AddCommand(describe.NewCmdDescribe("undistro", f, ioStreams))
	cmd.AddCommand(logs.NewCmdLogs(f, ioStreams))
	cmd.AddCommand(rollout.NewCmdRollout(f, ioStreams))
	cmd.AddCommand(apiresources.NewCmdAPIVersions(f, ioStreams))
	cmd.AddCommand(apiresources.NewCmdAPIResources(f, ioStreams))
	cmd.AddCommand(pinnipedcmd.LoginCmd)
	cmd.AddCommand(pinnipedcmd.NewWhoamiCommand(pinnipedcmd.GetRealConciergeClientset))
	cmd.AddCommand(NewCmdGet(f, ioStreams))
	cmd.AddCommand(NewCmdCreate(f, ioStreams))
	cmd.AddCommand(NewCmdInstall(cfgFlags, ioStreams))
	cmd.AddCommand(NewCmdMove(cfgFlags, ioStreams))
	cmd.AddCommand(NewCmdShowProgress(f, ioStreams))
	cmd.AddCommand(NewCmdUpgrade(f, ioStreams))
	cmd.AddCommand(NewCmdCompletion(ioStreams))
	cmd.AddCommand(version.NewVersionCommand())
	return cmd
}
