// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package shell

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/sh"
)

// Cmd represents a shell command
type Cmd struct {
	cmd  string
	args []string
}

// NewCmd shell command from a string
func NewCmd(cmd string) *Cmd {
	split := strings.Split(cmd, " ")
	return &Cmd{
		cmd:  split[0],
		args: split[1:],
	}
}

// NewCmdf creates a shell command from a format string
func NewCmdf(cmdf string, args ...interface{}) *Cmd {
	return NewCmd(fmt.Sprintf(cmdf, args...))
}

// Run runs shell command
func (c *Cmd) Run() error {
	return sh.Run(c.cmd, c.args...)
}

// RunWith runs shell command with environment variables
func (c *Cmd) RunWith(env map[string]string) error {
	return sh.RunWith(env, c.cmd, c.args...)
}

// Output runs the command and returns text from stdout
func (c *Cmd) Output() (string, error) {
	return sh.Output(c.cmd, c.args...)
}
