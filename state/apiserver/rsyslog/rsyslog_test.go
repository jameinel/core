// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package rsyslog_test

import (
	"encoding/pem"

	jc "github.com/juju/testing/checkers"
	gc "launchpad.net/gocheck"

	"github.com/wallyworld/core/juju/testing"
	"github.com/wallyworld/core/state"
	"github.com/wallyworld/core/state/api/params"
	apirsyslog "github.com/wallyworld/core/state/api/rsyslog"
	"github.com/wallyworld/core/state/apiserver/common"
	commontesting "github.com/wallyworld/core/state/apiserver/common/testing"
	"github.com/wallyworld/core/state/apiserver/rsyslog"
	apiservertesting "github.com/wallyworld/core/state/apiserver/testing"
	coretesting "github.com/wallyworld/core/testing"
)

type rsyslogSuite struct {
	testing.JujuConnSuite
	*commontesting.EnvironWatcherTest
	authorizer apiservertesting.FakeAuthorizer
	resources  *common.Resources
}

var _ = gc.Suite(&rsyslogSuite{})

func (s *rsyslogSuite) SetUpTest(c *gc.C) {
	s.JujuConnSuite.SetUpTest(c)
	s.authorizer = apiservertesting.FakeAuthorizer{
		LoggedIn:       true,
		EnvironManager: true,
	}
	s.resources = common.NewResources()
	api, err := rsyslog.NewRsyslogAPI(s.State, s.resources, s.authorizer)
	c.Assert(err, gc.IsNil)
	s.EnvironWatcherTest = commontesting.NewEnvironWatcherTest(
		api, s.State, s.resources, commontesting.NoSecrets)
}

func verifyRsyslogCACert(c *gc.C, st *apirsyslog.State, expected string) {
	cfg, err := st.EnvironConfig()
	c.Assert(err, gc.IsNil)
	c.Assert(cfg.RsyslogCACert(), gc.DeepEquals, expected)
}

func (s *rsyslogSuite) TestSetRsyslogCert(c *gc.C) {
	st, _ := s.OpenAPIAsNewMachine(c, state.JobManageEnviron)
	err := st.Rsyslog().SetRsyslogCert(coretesting.CACert)
	c.Assert(err, gc.IsNil)
	verifyRsyslogCACert(c, st.Rsyslog(), coretesting.CACert)
}

func (s *rsyslogSuite) TestSetRsyslogCertNil(c *gc.C) {
	st, _ := s.OpenAPIAsNewMachine(c, state.JobManageEnviron)
	err := st.Rsyslog().SetRsyslogCert("")
	c.Assert(err, gc.ErrorMatches, "no certificates found")
	verifyRsyslogCACert(c, st.Rsyslog(), "")
}

func (s *rsyslogSuite) TestSetRsyslogCertInvalid(c *gc.C) {
	st, _ := s.OpenAPIAsNewMachine(c, state.JobManageEnviron)
	err := st.Rsyslog().SetRsyslogCert(string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: []byte("not a valid certificate"),
	})))
	c.Assert(err, gc.ErrorMatches, ".*structure error.*")
	verifyRsyslogCACert(c, st.Rsyslog(), "")
}

func (s *rsyslogSuite) TestSetRsyslogCertPerms(c *gc.C) {
	st, _ := s.OpenAPIAsNewMachine(c, state.JobHostUnits)
	err := st.Rsyslog().SetRsyslogCert(coretesting.CACert)
	c.Assert(err, gc.ErrorMatches, "invalid entity name or password")
	c.Assert(err, jc.Satisfies, params.IsCodeUnauthorized)
	// Verify no change was effected.
	verifyRsyslogCACert(c, st.Rsyslog(), "")
}
