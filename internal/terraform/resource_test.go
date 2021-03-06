/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package terraform

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/terraform-docs/terraform-docs/internal/types"
)

func TestResourceFullType(t *testing.T) {
	assert := assert.New(t)
	resource := Resource{
		Type:           "private_key",
		ProviderName:   "tls",
		ProviderSource: "hashicorp/tls",
		Mode:           "managed",
		Version:        types.String("latest"),
	}
	assert.Equal("tls_private_key", resource.FullType())
}

func TestResourceURL(t *testing.T) {
	tests := []struct {
		name          string
		resource      Resource
		expectValue   string
		expectDefault bool
	}{
		{
			name: "generic URL construction",
			resource: Resource{
				Type:           "private_key",
				ProviderName:   "tls",
				ProviderSource: "hashicorp/tls",
				Mode:           "managed",
				Version:        types.String("latest"),
			},
			expectValue: "https://registry.terraform.io/providers/hashicorp/tls/latest/docs/resources/private_key",
		},
		{
			name: "unable to construct URL",
			resource: Resource{
				Type:           "custom",
				ProviderName:   "nih",
				ProviderSource: "http://nih.tld/some/path/to/provider/source",
				Mode:           "managed",
				Version:        types.String("latest"),
			},
			expectValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			assert.Equal(tt.expectValue, tt.resource.URL())
		})
	}
}

func TestResourcesSortedByType(t *testing.T) {
	assert := assert.New(t)
	resources := sampleResources()

	sort.Sort(resourcesSortedByType(resources))

	expected := []string{"a_a", "a_f", "b_b", "b_d", "c_c", "c_e", "c_e_x", "z_z", "z_z", "z_z"}
	actual := make([]string, len(resources))

	for k, i := range resources {
		actual[k] = i.FullType()
	}

	assert.Equal(expected, actual)
}

func TestResourcesSortedByTypeAndMode(t *testing.T) {
	assert := assert.New(t)
	resources := sampleResources()

	sort.Sort(resourcesSortedByType(resources))

	expected := []string{"a_a_m", "a_f_m", "b_b_m", "b_d_m", "c_c_m", "c_e_m", "c_e_x_m", "z_z", "z_z_d", "z_z_m"}
	actual := make([]string, len(resources))

	for k, i := range resources {
		v := i.FullType()
		switch i.Mode {
		case "managed":
			v = v + "_m"
		case "data":
			v = v + "_d"
		}
		actual[k] = v
	}

	assert.Equal(expected, actual)
}

func TestResourceVersion(t *testing.T) {
	tests := []struct {
		name       string
		constraint []string
		expected   string
	}{
		{
			name:       "exact version, without operator",
			constraint: []string{"1.2.3"},
			expected:   "1.2.3",
		},
		{
			name:       "exact version, with operator",
			constraint: []string{"= 1.2.3"},
			expected:   "1.2.3",
		},
		{
			name:       "exact version, with operator, without space",
			constraint: []string{"=1.2.3"},
			expected:   "1.2.3",
		},
		{
			name:       "exclude exact version, with space",
			constraint: []string{"!= 1.2.3"},
			expected:   "latest",
		},
		{
			name:       "exclude exact version, without space",
			constraint: []string{"!=1.2.3"},
			expected:   "latest",
		},
		{
			name:       "comparison version, with space",
			constraint: []string{"> 1.2.3"},
			expected:   "latest",
		},
		{
			name:       "comparison version, without space",
			constraint: []string{">1.2.3"},
			expected:   "latest",
		},
		{
			name:       "range version",
			constraint: []string{"> 1.2.3, < 2.0.0"},
			expected:   "latest",
		},
		{
			name:       "pessimistic version, with space",
			constraint: []string{"~> 1.2.3"},
			expected:   "latest",
		},
		{
			name:       "pessimistic version, without space",
			constraint: []string{"~>1.2.3"},
			expected:   "latest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			assert.Equal(tt.expected, resourceVersion(tt.constraint))
		})
	}
}

func sampleResources() []*Resource {
	return []*Resource{
		{
			Type:           "e",
			ProviderName:   "c",
			ProviderSource: "hashicorp/e",
			Mode:           "managed",
			Version:        "1.5.0",
		},
		{
			Type:           "e_x",
			ProviderName:   "c",
			ProviderSource: "hashicorp/e",
			Mode:           "managed",
			Version:        "1.5.0",
		},
		{
			Type:           "a",
			ProviderName:   "a",
			ProviderSource: "hashicorp/a",
			Mode:           "managed",
			Version:        "1.1.0",
		},
		{
			Type:           "d",
			ProviderName:   "b",
			ProviderSource: "hashicorp/d",
			Mode:           "managed",
			Version:        "1.4.0",
		},
		{
			Type:           "b",
			ProviderName:   "b",
			ProviderSource: "hashicorp/b",
			Mode:           "managed",
			Version:        "1.2.0",
		},
		{
			Type:           "c",
			ProviderName:   "c",
			ProviderSource: "hashicorp/c",
			Mode:           "managed",
			Version:        "1.3.0",
		},
		{
			Type:           "f",
			ProviderName:   "a",
			ProviderSource: "hashicorp/f",
			Mode:           "managed",
			Version:        "1.6.0",
		},
		{
			ProviderName:   "z",
			Type:           "z",
			ProviderSource: "hashicorp/a",
			Mode:           "managed",
			Version:        "1.5.0",
		},
		{
			ProviderName:   "z",
			Type:           "z",
			ProviderSource: "hashicorp/a",
			Mode:           "data",
			Version:        "1.5.0",
		},
		{
			ProviderName:   "z",
			Type:           "z",
			ProviderSource: "hashicorp/a",
			Mode:           "",
			Version:        "1.5.0",
		},
	}
}
