package juju_test

import (
	"github.com/wallyworld/core/instance"
	"github.com/wallyworld/core/juju"
	"github.com/wallyworld/core/state/api"
)

type mockAPIState struct {
	close func(juju.APIState) error

	apiHostPorts [][]instance.HostPort
}

func (s *mockAPIState) Close() error {
	if s.close != nil {
		return s.close(s)
	}
	return nil
}

func (s *mockAPIState) APIHostPorts() [][]instance.HostPort {
	return s.apiHostPorts
}

func panicAPIOpen(apiInfo *api.Info, opts api.DialOpts) (juju.APIState, error) {
	panic("api.Open called unexpectedly")
}
