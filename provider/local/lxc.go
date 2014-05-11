// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package local

import (
	"strconv"

	"github.com/wallyworld/core/instance"
	"github.com/wallyworld/core/version"
)

// releaseVersion is a function that returns a string representing the
// DISTRIB_RELEASE from the /etc/lsb-release file.
var releaseVersion = version.ReleaseVersion

func useFastLXC(containerType instance.ContainerType) bool {
	if containerType != instance.LXC {
		return false
	}
	release := releaseVersion()
	if release == "" {
		return false
	}
	value, err := strconv.ParseFloat(release, 64)
	if err != nil {
		return false
	}
	return value >= 14.04
}
