package osenv_test

import (
	gc "launchpad.net/gocheck"

	"github.com/wallyworld/core/juju/osenv"
)

func (s *importSuite) TestHomeLinux(c *gc.C) {
	h := "/home/foo/bar"
	s.PatchEnvironment("HOME", h)
	c.Check(osenv.Home(), gc.Equals, h)
}
