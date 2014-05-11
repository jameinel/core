// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common

import (
	"github.com/wallyworld/core/environs"
	"github.com/wallyworld/core/environs/imagemetadata"
	"github.com/wallyworld/core/environs/simplestreams"
	"github.com/wallyworld/core/utils/set"
)

// SupportedArchitectures returns all the image architectures for env matching the constraints.
func SupportedArchitectures(env environs.Environ, imageConstraint *imagemetadata.ImageConstraint) ([]string, error) {
	sources, err := imagemetadata.GetMetadataSources(env)
	if err != nil {
		return nil, err
	}
	matchingImages, _, err := imagemetadata.Fetch(sources, simplestreams.DefaultIndexPath, imageConstraint, false)
	if err != nil {
		return nil, err
	}
	var arches = set.NewStrings()
	for _, im := range matchingImages {
		arches.Add(im.Arch)
	}
	return arches.Values(), nil
}
