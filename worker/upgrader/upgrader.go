// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upgrader

import (
	"fmt"
	"net/http"
	"time"

	"github.com/loggo/loggo"
	"launchpad.net/tomb"

	"launchpad.net/juju-core/agent"
	agenttools "launchpad.net/juju-core/agent/tools"
	apiwatcher "launchpad.net/juju-core/state/api/watcher"
	"launchpad.net/juju-core/state/watcher"
	coretools "launchpad.net/juju-core/tools"
	"launchpad.net/juju-core/utils"
	"launchpad.net/juju-core/version"
)

// retryAfter returns a channel that receives a value
// when a failed download should be retried.
var retryAfter = func() <-chan time.Time {
	return time.After(5 * time.Second)
}

var logger = loggo.GetLogger("juju.worker.upgrader")

// upgraderClient is implemented by whatever client facade is used by this upgrade worker.
type upgradeClient interface {
	DesiredVersion(tag string) (version.Number, error)
	WatchAPIVersion(agentTag string) (apiwatcher.NotifyWatcher, error)
	SetVersion(tag string, v version.Binary) error
	Tools(tag string) (*coretools.Tools, bool, error)
}

// Upgrader represents a worker that watches the state for upgrade
// requests.
type Upgrader struct {
	tomb    tomb.Tomb
	client  upgradeClient
	dataDir string
	tag     string
}

// NewUpgrader returns a new upgrader worker. It watches changes to the
// current version of the current agent (with the given tag) and tries to
// download the tools for any new version into the given data directory.  If
// an upgrade is needed, the worker will exit with an UpgradeReadyError
// holding details of the requested upgrade. The tools will have been
// downloaded and unpacked.
func NewUpgrader(client upgradeClient, agentConfig agent.Config) *Upgrader {
	u := &Upgrader{
		client:  client,
		dataDir: agentConfig.DataDir(),
		tag:     agentConfig.Tag(),
	}
	go func() {
		defer u.tomb.Done()
		u.tomb.Kill(u.loop())
	}()
	return u
}

// Kill implements worker.Worker.Kill.
func (u *Upgrader) Kill() {
	u.tomb.Kill(nil)
}

// Wait implements worker.Worker.Wait.
func (u *Upgrader) Wait() error {
	return u.tomb.Wait()
}

// Stop stops the upgrader and returns any
// error it encountered when running.
func (u *Upgrader) Stop() error {
	u.Kill()
	return u.Wait()
}

func (u *Upgrader) loop() error {
	currentTools := &coretools.Tools{Version: version.Current}
	err := u.client.SetVersion(u.tag, currentTools.Version)
	if err != nil {
		return err
	}
	versionWatcher, err := u.client.WatchAPIVersion(u.tag)
	if err != nil {
		return err
	}
	changes := versionWatcher.Changes()
	defer watcher.Stop(versionWatcher, &u.tomb)
	var retry <-chan time.Time
	// We don't read on the dying channel until we have received the
	// initial event from the API version watcher, thus ensuring
	// that we attempt an upgrade even if other workers are dying
	// all around us.
	var (
		dying                          <-chan struct{}
		wantTools                      *coretools.Tools
		wantVersion                    version.Number
		disableSSLHostnameVerification bool
	)
	for {
		select {
		case _, ok := <-changes:
			if !ok {
				return watcher.MustErr(versionWatcher)
			}
			wantVersion, err = u.client.DesiredVersion(u.tag)
			if err != nil {
				return err
			}
			logger.Infof("desired tool version: %v", wantVersion)
			dying = u.tomb.Dying()
		case <-retry:
		case <-dying:
			return nil
		}
		if wantVersion != currentTools.Version.Number {
			logger.Infof("upgrade requested from %v to %v", currentTools.Version, wantVersion)
			// TODO(dimitern) 2013-10-03 bug #1234715
			// Add a testing HTTPS storage to verify the
			// disableSSLHostnameVerification behavior here.
			wantTools, disableSSLHostnameVerification, err = u.client.Tools(u.tag)
			if err != nil {
				// Not being able to lookup Tools is considered fatal
				return err
			}
			// The worker cannot be stopped while we're downloading
			// the tools - this means that even if the API is going down
			// repeatedly (causing the agent to be stopped), as long
			// as we have got as far as this, we will still be able to
			// upgrade the agent.
			err := u.ensureTools(wantTools, disableSSLHostnameVerification)
			if err == nil {
				return &UpgradeReadyError{
					OldTools:  version.Current,
					NewTools:  wantTools.Version,
					AgentName: u.tag,
					DataDir:   u.dataDir,
				}
			}
			logger.Errorf("failed to fetch tools from %q: %v", wantTools.URL, err)
			retry = retryAfter()
		}
	}
}

func (u *Upgrader) ensureTools(agentTools *coretools.Tools, disableSSLHostnameVerification bool) error {
	if _, err := agenttools.ReadTools(u.dataDir, agentTools.Version); err == nil {
		// Tools have already been downloaded
		return nil
	}
	client := http.DefaultClient
	logger.Infof("fetching tools from %q", agentTools.URL)
	if disableSSLHostnameVerification {
		logger.Infof("hostname SSL verification disabled")
		client = utils.GetNonValidatingHTTPClient()
	}
	resp, err := client.Get(agentTools.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP response: %v", resp.Status)
	}
	err = agenttools.UnpackTools(u.dataDir, agentTools, resp.Body)
	if err != nil {
		return fmt.Errorf("cannot unpack tools: %v", err)
	}
	return nil
}
