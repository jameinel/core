// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package tools

import (
	"fmt"

	"launchpad.net/juju-core/constraints"
	"launchpad.net/juju-core/environs"
	"launchpad.net/juju-core/errors"
	"launchpad.net/juju-core/log"
	coretools "launchpad.net/juju-core/tools"
	"launchpad.net/juju-core/version"
	"launchpad.net/loggo"
)

var logger = loggo.GetLogger("juju.environs.tools")

// FindTools returns a List containing all tools with a given
// major version number available in the environment, filtered by filter.
// If *any* tools are present in private storage, *only* tools from private
// storage are available.
// If *no* tools are present in private storage, *only* tools from public
// storage are available.
// If no *available* tools have the supplied major version number, or match the
// supplied filter, the function returns a *NotFoundError.
func FindTools(environ environs.Environ, majorVersion int, filter coretools.Filter) (list coretools.List, err error) {
	log.Infof("environs: reading tools with major version %d", majorVersion)
	defer convertToolsError(&err)
	// Construct a tools filter.
	// Discard all that are known to be irrelevant.
	if filter.Number != version.Zero {
		log.Infof("environs: filtering tools by version: %s", filter.Number.Major)
	}
	if filter.Series != "" {
		log.Infof("environs: filtering tools by series: %s", filter.Series)
	}
	if filter.Arch != "" {
		log.Infof("environs: filtering tools by architecture: %s", filter.Arch)
	}
	list, err = ReadList(environ.Storage(), majorVersion)
	if err == ErrNoTools {
		log.Infof("environs: falling back to public bucket")
		list, err = ReadList(environ.PublicStorage(), majorVersion)
	}
	if err != nil {
		return nil, err
	}
	return list.Match(filter)
}

// FindBootstrapTools returns a ToolsList containing only those tools with
// which it would be reasonable to launch an environment's first machine,
// given the supplied constraints.
// If the environment was not already configured to use a specific agent
// version, the newest available version will be chosen and set in the
// environment's configuration.
func FindBootstrapTools(environ environs.Environ, cons constraints.Value) (list coretools.List, err error) {
	// Construct a tools filter.
	cliVersion := version.Current.Number
	cfg := environ.Config()
	filter := coretools.Filter{
		Series: cfg.DefaultSeries(),
		Arch:   stringOrEmpty(cons.Arch),
	}
	if agentVersion, ok := cfg.AgentVersion(); ok {
		// If we already have an explicit agent version set, we're done.
		filter.Number = agentVersion
		return FindTools(environ, cliVersion.Major, filter)
	}
	if dev := cliVersion.IsDev() || cfg.Development(); !dev {
		log.Infof("environs: filtering tools by released version")
		filter.Released = true
	}
	list, err = FindTools(environ, cliVersion.Major, filter)
	if err != nil {
		return nil, err
	}

	// We probably still have a mix of versions available; discard older ones
	// and update environment configuration to use only those remaining.
	agentVersion, list := list.Newest()
	log.Infof("environs: picked newest version: %s", agentVersion)
	cfg, err = cfg.Apply(map[string]interface{}{
		"agent-version": agentVersion.String(),
	})
	if err == nil {
		err = environ.SetConfig(cfg)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update environment configuration: %v", err)
	}
	return list, nil
}

func stringOrEmpty(pstr *string) string {
	if pstr == nil {
		return ""
	}
	return *pstr
}

// FindInstanceTools returns a ToolsList containing only those tools with which
// it would be reasonable to start a new instance, given the supplied series and
// constraints.
// It is an error to call it with an environment not already configured to use
// a specific agent version.
func FindInstanceTools(environ environs.Environ, series string, cons constraints.Value) (list coretools.List, err error) {
	// Construct a tools filter.
	// Discard all that are known to be irrelevant.
	agentVersion, ok := environ.Config().AgentVersion()
	if !ok {
		return nil, fmt.Errorf("no agent version set in environment configuration")
	}
	filter := coretools.Filter{
		Number: agentVersion,
		Series: series,
		Arch:   stringOrEmpty(cons.Arch),
	}
	return FindTools(environ, agentVersion.Major, filter)
}

// FindExactTools returns only the tools that match the supplied version.
// TODO(fwereade) this should not exist: it's used by cmd/jujud/Upgrader,
// which needs to run on every agent and must absolutely *not* in general
// have access to an environs.Environ.
func FindExactTools(environ environs.Environ, vers version.Binary) (t *coretools.Tools, err error) {
	log.Infof("environs: finding exact version %s", vers)
	filter := coretools.Filter{
		Number: vers.Number,
		Series: vers.Series,
		Arch:   vers.Arch,
	}
	list, err := FindTools(environ, vers.Number.Major, filter)
	if err != nil {
		return nil, err
	}
	return list[0], nil
}

// CheckToolsSeries verifies that all the given possible tools are for the
// given OS series.
func CheckToolsSeries(toolsList coretools.List, series string) error {
	toolsSeries := toolsList.Series()
	if len(toolsSeries) != 1 {
		return fmt.Errorf("expected single series, got %v", toolsSeries)
	}
	if toolsSeries[0] != series {
		return fmt.Errorf("tools mismatch: expected series %v, got %v", series, toolsSeries[0])
	}
	return nil
}

func isToolsError(err error) bool {
	switch err {
	case ErrNoTools, coretools.ErrNoMatches:
		return true
	}
	return false
}

func convertToolsError(err *error) {
	if isToolsError(*err) {
		*err = errors.NewNotFoundError(*err, "")
	}
}