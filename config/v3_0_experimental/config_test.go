// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v3_0_experimental

import (
	"testing"

	"github.com/coreos/ignition/config/shared/errors"
	"github.com/coreos/ignition/config/v3_0_experimental/types"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type in struct {
		config []byte
	}
	type out struct {
		config         types.Config
		err            error
		checkOnStrings bool
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{config: []byte(`{"ignitionVersion": 1}`)},
			out: out{config: types.Config{Ignition: types.Ignition{Version: types.MaxVersion.String()}}},
		},
		{
			in:  in{config: []byte(`{"ignition": {"version": "1.0.0"}}`)},
			out: out{err: errors.ErrUnknownVersion},
		},
		{
			in:  in{config: []byte(`{"ignition": {"version": "2.0.0"}}`)},
			out: out{config: types.Config{Ignition: types.Ignition{Version: types.MaxVersion.String()}}},
		},
		{
			in:  in{config: []byte(`{"ignition": {"version": "2.1.0"}}`)},
			out: out{config: types.Config{Ignition: types.Ignition{Version: types.MaxVersion.String()}}},
		},
		{
			in:  in{config: []byte(`{"ignition": {"version": "2.2.0-experimental"}}`)},
			out: out{err: errors.ErrUnknownVersion},
		},
		{
			in:  in{config: []byte(`{"ignition": {"version": "2.1.0-experimental"}}`)},
			out: out{err: errors.ErrUnknownVersion},
		},
		{
			in:  in{config: []byte(`{"ignition": {"version": "2.3.0-experimental"}}`)},
			out: out{err: errors.ErrUnknownVersion},
		},
		{
			in:  in{config: []byte(`{"ignition": {"version": "2.2.0"}}`)},
			out: out{config: types.Config{Ignition: types.Ignition{Version: types.MaxVersion.String()}}},
		},
		{
			in:  in{config: []byte(`{"ignition": {"version": "2.3.0"}}`)},
			out: out{config: types.Config{Ignition: types.Ignition{Version: types.MaxVersion.String()}}},
		},
		{
			in:  in{config: []byte(`{"ignition": {"version": "2.4.0-experimental"}}`)},
			out: out{config: types.Config{Ignition: types.Ignition{Version: types.MaxVersion.String()}}},
		},
		{
			in:  in{config: []byte(`{"ignition": {"version": "2.0.0"},}`)},
			out: out{err: errors.ErrInvalid},
		},
		{
			in:  in{config: []byte(`{"ignition": {"version": "invalid.semver"}}`)},
			out: out{err: errors.ErrUnknownVersion},
		},
		{
			in:  in{config: []byte(`{}`)},
			out: out{err: errors.ErrUnknownVersion},
		},
		{
			in:  in{config: []byte{}},
			out: out{err: errors.ErrEmpty},
		},
		{
			in:  in{config: []byte("#cloud-config")},
			out: out{err: errors.ErrCloudConfig},
		},
		{
			in:  in{config: []byte("#cloud-config ")},
			out: out{err: errors.ErrCloudConfig},
		},
		{
			in:  in{config: []byte("#cloud-config\n\r")},
			out: out{err: errors.ErrCloudConfig},
		},
		{
			in: in{config: []byte{0x1f, 0x8b, 0x08, 0x00, 0x03, 0xd6, 0x79, 0x56,
				0x00, 0x03, 0x53, 0x4e, 0xce, 0xc9, 0x2f, 0x4d, 0xd1, 0x4d, 0xce,
				0xcf, 0x4b, 0xcb, 0x4c, 0xe7, 0x02, 0x00, 0x05, 0x56, 0xb3, 0xb8,
				0x0e, 0x00, 0x00, 0x00}},
			out: out{err: errors.ErrCloudConfig},
		},
		{
			in:  in{config: []byte("#!/bin/sh")},
			out: out{err: errors.ErrScript},
		},
		{
			in: in{config: []byte{0x1f, 0x8b, 0x08, 0x00, 0x48, 0xda, 0x79, 0x56,
				0x00, 0x03, 0x53, 0x56, 0xd4, 0x4f, 0xca, 0xcc, 0xd3, 0x2f, 0xce,
				0xe0, 0x02, 0x00, 0x1d, 0x9d, 0xfb, 0x04, 0x0a, 0x00, 0x00, 0x00}},
			out: out{err: errors.ErrScript},
		},
	}

	for i, test := range tests {
		config, report, err := Parse(test.in.config)
		if (!test.out.checkOnStrings && test.out.err != err) ||
			(test.out.checkOnStrings && test.out.err.Error() != err.Error()) {
			t.Errorf("#%d: bad error: want %v, got %v, report: %+v", i, test.out.err, err, report)
		}
		assert.Equal(t, test.out.config, config, "#%d: bad config, report: %+v", i, report)
	}
}
