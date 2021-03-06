// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"fmt"

	"github.com/wallyworld/core/cmd"
	"github.com/wallyworld/core/cmd/envcmd"
	"github.com/wallyworld/core/juju"
)

// DestroyRelationCommand causes an existing service relation to be shut down.
type DestroyRelationCommand struct {
	envcmd.EnvCommandBase
	Endpoints []string
}

func (c *DestroyRelationCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "destroy-relation",
		Args:    "<service1>[:<relation name1>] <service2>[:<relation name2>]",
		Purpose: "destroy a relation between two services",
		Aliases: []string{"remove-relation"},
	}
}

func (c *DestroyRelationCommand) Init(args []string) error {
	if err := c.EnsureEnvName(); err != nil {
		return err
	}
	if len(args) != 2 {
		return fmt.Errorf("a relation must involve two services")
	}
	c.Endpoints = args
	return nil
}

func (c *DestroyRelationCommand) Run(_ *cmd.Context) error {
	client, err := juju.NewAPIClientFromName(c.EnvName)
	if err != nil {
		return err
	}
	defer client.Close()
	return client.DestroyRelation(c.Endpoints...)
}
