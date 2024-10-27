// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
)

type Option func(_ *index) error

// WithDiscoverer returns Option that sets discoverer client.
func WithDiscoverer(client discoverer.Client) Option {
	return func(idx *index) error {
		if client == nil {
			return errors.NewErrCriticalOption("discoverer", client)
		}
		idx.client = client
		return nil
	}
}

// WithIndexingConcurrency returns Option that sets indexing concurrency.
func WithIndexingConcurrency(num int) Option {
	return func(idx *index) error {
		if num <= 0 {
			return errors.NewErrInvalidOption("indexingConcurrency", num)
		}
		idx.concurrency = num
		return nil
	}
}

// WithTargetAddrs returns Option that sets indexing target addresses.
func WithTargetAddrs(addrs ...string) Option {
	return func(idx *index) error {
		if len(addrs) != 0 {
			idx.targetAddrs = append(idx.targetAddrs, addrs...)
		}
		return nil
	}
}

// WithTargetIndexID returns Option that sets target deleting index ID.
func WithTargetIndexID(indexID string) Option {
	return func(idx *index) error {
		idx.targetIndexID = indexID
		return nil
	}
}
