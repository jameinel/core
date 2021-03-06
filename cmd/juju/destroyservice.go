// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"fmt"

	"github.com/wallyworld/core/cmd"
	"github.com/wallyworld/core/cmd/envcmd"
	"github.com/wallyworld/core/juju"
	"github.com/wallyworld/core/names"
)

// DestroyServiceCommand causes an existing service to be destroyed.
type DestroyServiceCommand struct {
	envcmd.EnvCommandBase
	ServiceName string
}

func (c *DestroyServiceCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "destroy-service",
		Args:    "<service>",
		Purpose: "destroy a service",
		Doc:     "Destroying a service will destroy all its units and relations.",
		Aliases: []string{"remove-service"},
	}
}

func (c *DestroyServiceCommand) Init(args []string) error {
	if err := c.EnsureEnvName(); err != nil {
		return err
	}
	if len(args) == 0 {
		return fmt.Errorf("no service specified")
	}
	if !names.IsService(args[0]) {
		return fmt.Errorf("invalid service name %q", args[0])
	}
	c.ServiceName, args = args[0], args[1:]
	return cmd.CheckEmpty(args)
}

func (c *DestroyServiceCommand) Run(_ *cmd.Context) error {
	client, err := juju.NewAPIClientFromName(c.EnvName)
	if err != nil {
		return err
	}
	defer client.Close()
	return client.ServiceDestroy(c.ServiceName)
}
