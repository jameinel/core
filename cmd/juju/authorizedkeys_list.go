// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"fmt"
	"strings"

	"launchpad.net/gnuflag"

	"github.com/wallyworld/core/cmd"
	"github.com/wallyworld/core/cmd/envcmd"
	"github.com/wallyworld/core/juju"
	"github.com/wallyworld/core/utils/ssh"
)

var listKeysDoc = `
List a user's authorized ssh keys, allowing the holders of those keys to log on to Juju nodes.
By default, just the key fingerprint is printed. Use --full to display the entire key.

`

// ListKeysCommand is used to list the authorized ssh keys.
type ListKeysCommand struct {
	envcmd.EnvCommandBase
	showFullKey bool
	user        string
}

func (c *ListKeysCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "list",
		Doc:     listKeysDoc,
		Purpose: "list authorized ssh keys for a specified user",
	}
}

func (c *ListKeysCommand) SetFlags(f *gnuflag.FlagSet) {
	c.EnvCommandBase.SetFlags(f)
	f.BoolVar(&c.showFullKey, "full", false, "show full key instead of just the key fingerprint")
	f.StringVar(&c.user, "user", "admin", "the user for which to list the keys")
}

func (c *ListKeysCommand) Init(args []string) error {
	if err := c.EnvCommandBase.EnsureEnvName(); err != nil {
		return err
	}
	return cmd.CheckEmpty(args)
}

func (c *ListKeysCommand) Run(context *cmd.Context) error {
	client, err := juju.NewKeyManagerClient(c.EnvName)
	if err != nil {
		return err
	}
	defer client.Close()

	mode := ssh.Fingerprints
	if c.showFullKey {
		mode = ssh.FullKeys
	}
	results, err := client.ListKeys(mode, c.user)
	if err != nil {
		return err
	}
	result := results[0]
	if result.Error != nil {
		return result.Error
	}
	fmt.Fprintf(context.Stdout, "Keys for user %s:\n", c.user)
	fmt.Fprintln(context.Stdout, strings.Join(result.Result, "\n"))
	return nil
}
