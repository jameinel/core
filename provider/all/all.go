// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package all

// Register all the available providers.
import (
	_ "github.com/wallyworld/core/provider/azure"
	_ "github.com/wallyworld/core/provider/ec2"
	_ "github.com/wallyworld/core/provider/joyent"
	_ "github.com/wallyworld/core/provider/local"
	_ "github.com/wallyworld/core/provider/maas"
	_ "github.com/wallyworld/core/provider/manual"
	_ "github.com/wallyworld/core/provider/openstack"
)
