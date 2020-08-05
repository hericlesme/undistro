/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package util

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Cmd implements a wrapper on os/exec.cmd
type Cmd struct {
	command string
	args    []string
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
}

func NewCmd(command string, args ...string) *Cmd {
	return &Cmd{
		command: command,
		args:    args,
	}
}

func (c *Cmd) Run() error {
	return c.runInnerCommand()
}

func (c *Cmd) RunWithEcho() error {
	c.stdout = os.Stderr
	c.stderr = os.Stdout
	return c.runInnerCommand()
}

func (c *Cmd) RunAndCapture() (lines []string, err error) {
	var buff bytes.Buffer
	c.stdout = &buff
	c.stderr = &buff
	err = c.runInnerCommand()

	scanner := bufio.NewScanner(&buff)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())

	}
	return lines, err
}

func (c *Cmd) Stdin(in io.Reader) *Cmd {
	c.stdin = in
	return c
}

func (c *Cmd) runInnerCommand() error {
	cmd := exec.Command(c.command, c.args...) //nolint:gosec

	if c.stdin != nil {
		cmd.Stdin = c.stdin
	}
	if c.stdout != nil {
		cmd.Stdout = c.stdout
	}
	var b1 strings.Builder
	cmd.Stderr = &b1
	if c.stderr != nil {
		cmd.Stderr = io.MultiWriter(&b1, c.stderr)
	}

	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to run: %s %s\n%s", c.command, strings.Join(c.args, " "), b1.String())
	}

	return nil
}
