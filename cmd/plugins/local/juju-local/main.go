// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"os"

	"github.com/wallyworld/core/cmd/plugins/local"

	// Import only the local provider.
	_ "github.com/wallyworld/core/provider/local"
)

func main() {
	local.Main(os.Args)
}
