// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"fmt"
	"os"

	"github.com/juju/loggo"

	"github.com/wallyworld/core/cmd"
	"github.com/wallyworld/core/juju"
	_ "github.com/wallyworld/core/provider/all"
)

var logger = loggo.GetLogger("juju.plugins.metadata")

var metadataDoc = `
Juju metadata is used to find the correct image and tools when bootstrapping a
Juju environment.
`

// Main registers subcommands for the juju-metadata executable, and hands over control
// to the cmd package. This function is not redundant with main, because it
// provides an entry point for testing with arbitrary command line arguments.
func Main(args []string) {
	ctx, err := cmd.DefaultContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	if err := juju.InitJujuHome(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}
	metadatacmd := cmd.NewSuperCommand(cmd.SuperCommandParams{
		Name:        "metadata",
		UsagePrefix: "juju",
		Doc:         metadataDoc,
		Purpose:     "tools for generating and validating image and tools metadata",
		Log:         &cmd.Log{}})

	metadatacmd.Register(&ValidateImageMetadataCommand{})
	metadatacmd.Register(&ImageMetadataCommand{})
	metadatacmd.Register(&ToolsMetadataCommand{})
	metadatacmd.Register(&ValidateToolsMetadataCommand{})
	metadatacmd.Register(&SignMetadataCommand{})

	os.Exit(cmd.Main(metadatacmd, ctx, args[1:]))
}

func main() {
	Main(os.Args)
}
