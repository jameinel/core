// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upgrader

import (
	"github.com/wallyworld/core/tools"
	"github.com/wallyworld/core/utils"
)

var (
	RetryAfter           = &retryAfter
	AllowedTargetVersion = allowedTargetVersion
)

func EnsureTools(u *Upgrader, agentTools *tools.Tools, hostnameVerification utils.SSLHostnameVerification) error {
	return u.ensureTools(agentTools, hostnameVerification)
}
