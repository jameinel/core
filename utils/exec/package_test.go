// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package exec_test

import (
	"testing"

	gc "launchpad.net/gocheck"

	"github.com/wallyworld/core/testing/testbase"
)

func Test(t *testing.T) { gc.TestingT(t) }

type Dependencies struct{}

var _ = gc.Suite(&Dependencies{})

func (*Dependencies) TestPackageDependencies(c *gc.C) {
	// This test is to ensure we don't bring in dependencies without thinking.
	c.Assert(testbase.FindJujuCoreImports(c, "github.com/wallyworld/core/utils/exec"),
		gc.HasLen, 0)
}
