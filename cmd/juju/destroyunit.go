// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"errors"
	"fmt"

	"github.com/wallyworld/core/cmd"
	"github.com/wallyworld/core/cmd/envcmd"
	"github.com/wallyworld/core/juju"
	"github.com/wallyworld/core/names"
)

// DestroyUnitCommand is responsible for destroying service units.
type DestroyUnitCommand struct {
	envcmd.EnvCommandBase
	UnitNames []string
}

func (c *DestroyUnitCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "destroy-unit",
		Args:    "<unit> [...]",
		Purpose: "destroy service units",
		Aliases: []string{"remove-unit"},
	}
}

func (c *DestroyUnitCommand) Init(args []string) error {
	if err := c.EnsureEnvName(); err != nil {
		return err
	}
	c.UnitNames = args
	if len(c.UnitNames) == 0 {
		return errors.New("no units specified")
	}
	for _, name := range c.UnitNames {
		if !names.IsUnit(name) {
			return fmt.Errorf("invalid unit name %q", name)
		}
	}
	return nil
}

// Run connects to the environment specified on the command line and destroys
// units therein.
func (c *DestroyUnitCommand) Run(_ *cmd.Context) error {
	client, err := juju.NewAPIClientFromName(c.EnvName)
	if err != nil {
		return err
	}
	defer client.Close()
	return client.DestroyServiceUnits(c.UnitNames...)
}
