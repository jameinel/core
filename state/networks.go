// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"fmt"

	"labix.org/v2/mgo/bson"

	"launchpad.net/juju-core/names"
)

// Network represents the state of a network.
type Network struct {
	st  *State
	doc networkDoc
}

// NetworkInfo describes a single network.
type NetworkInfo struct {
	// Name is juju-internal name of the network.
	Name string

	// ProviderId is a provider-specific network id.
	ProviderId string

	// CIDR of the network, in 123.45.67.89/24 format.
	CIDR string

	// VLANTag needs to be between 1 and 4094 for VLANs and 0 for
	// normal networks. It's defined by IEEE 802.1Q standard.
	VLANTag int

	// IsVirtual is true when this network uses virtual network
	// interface devices, or false when using physical devices.
	IsVirtual bool
}

// networkDoc represents a configured network that a machine can be a
// part of.
type networkDoc struct {
	// Name is the network's name. It should be one of the machine's
	// included networks.
	Name string `bson:"_id"`

	ProviderId string
	CIDR       string
	VLANTag    int
	IsVirtual  bool
}

func newNetwork(st *State, doc *networkDoc) *Network {
	return &Network{st, *doc}
}

func newNetworkDoc(args NetworkInfo) *networkDoc {
	return &networkDoc{
		Name:       args.Name,
		ProviderId: args.ProviderId,
		CIDR:       args.CIDR,
		VLANTag:    args.VLANTag,
		IsVirtual:  args.IsVirtual,
	}
}

// GoString implements fmt.GoStringer.
func (n *Network) GoString() string {
	return fmt.Sprintf(
		"&state.Network{name: %q, providerId: %q, cidr: %q, vlanTag: %v, isVirtual: %v}",
		n.Name(), n.ProviderId(), n.CIDR(), n.VLANTag(), n.IsVirtual())
}

// Name returns the network name.
func (n *Network) Name() string {
	return n.doc.Name
}

// ProviderId returns the provider-specific id of the network.
func (n *Network) ProviderId() string {
	return n.doc.ProviderId
}

// Tag returns the network tag.
func (n *Network) Tag() string {
	return names.NetworkTag(n.doc.Name)
}

// CIDR returns the network CIDR (e.g. 192.168.50.0/24).
func (n *Network) CIDR() string {
	return n.doc.CIDR
}

// VLANTag returns the network VLAN tag. It's a number between 1 and
// 4094 for VLANs and 0 if the network is not a VLAN.
func (n *Network) VLANTag() int {
	return n.doc.VLANTag
}

// IsVLAN returns whether the network is a VLAN (has tag > 0) or a
// normal network.
func (n *Network) IsVLAN() bool {
	return n.doc.VLANTag > 0
}

// IsVirtual returns whether the network is virtual network (using
// virtual network interfaces only).
func (n *Network) IsVirtual() bool {
	return n.doc.IsVirtual
}

// IsPhysical returns whether the network is physical network (using
// physical network interfaces only).
func (n *Network) IsPhysical() bool {
	return !n.doc.IsVirtual
}

// Interfaces returns all network interfaces on the network.
func (n *Network) Interfaces() ([]*NetworkInterface, error) {
	docs := []networkInterfaceDoc{}
	sel := bson.D{{"networkname", n.doc.Name}}
	err := n.st.networkInterfaces.Find(sel).All(&docs)
	if err != nil {
		return nil, err
	}
	ifaces := make([]*NetworkInterface, len(docs))
	for i, doc := range docs {
		ifaces[i] = newNetworkInterface(n.st, &doc)
	}
	return ifaces, nil
}
