// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upgrades

import (
	"github.com/wallyworld/core/environs/config"
)

func updateRsyslogPort(context Context) error {
	st := context.State()
	attrs := map[string]interface{}{
		"syslog-port": config.DefaultSyslogPort,
	}
	return st.UpdateEnvironConfig(attrs, nil, nil)
}
