// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package bbolt_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/db/kvs/bbolt"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

const mode = os.FileMode(0o600)

func TestNew(t *testing.T) {
	t.Parallel()

	type test struct {
		name     string
		testfunc func(t *testing.T)
	}

	tests := []test{
		{
			name: "New returns bbolt instance with new file when file does not exist",
			testfunc: func(t *testing.T) {
				tempdir := t.TempDir()
				tmpfile := filepath.Join(tempdir, "test.db")

				b, err := bbolt.New(tmpfile, "", mode)
				require.NoError(t, err)
				require.NotNil(t, b)
			},
		},
		{
			name: "New returns bbolt instance with existing file",
			testfunc: func(t *testing.T) {
				tempdir := t.TempDir()
				tmpfile := filepath.Join(tempdir, "test.db")

				// create a file
				f, err := os.Create(tmpfile)
				require.NoError(t, err)
				err = f.Close()
				require.NoError(t, err)

				b, err := bbolt.New(f.Name(), "", mode)
				require.NoError(t, err)
				require.NotNil(t, b)
			},
		},
		{
			name: "New returns bbolt with custom bucket name",
			testfunc: func(t *testing.T) {
				tempdir := t.TempDir()
				tmpfile := filepath.Join(tempdir, "test.db")

				b, err := bbolt.New(tmpfile, "my bucket name", mode)
				require.NoError(t, err)
				require.NotNil(t, b)
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			test.testfunc(tt)
		})
	}
}

func Test_bbolt_GetSetClose(t *testing.T) {
	t.Parallel()

	type test struct {
		name     string
		testfunc func(t *testing.T)
	}

	setup := func(t *testing.T) (b bbolt.Bbolt, file string) {
		tempdir := t.TempDir()
		tmpfile := filepath.Join(tempdir, "test.db")
		b, err := bbolt.New(tmpfile, "", mode)
		require.NoError(t, err)

		return b, tmpfile
	}

	tests := []test{
		{
			name: "Succeed to set and get with the key returns the value",
			testfunc: func(t *testing.T) {
				b, _ := setup(t)

				k, v := []byte("key"), []byte("value")
				err := b.Set(k, v)
				require.NoError(t, err)

				val, ok, err := b.Get(k)
				require.NoError(t, err)
				require.True(t, ok)
				require.Equal(t, v, val)
			},
		},
		{
			name: "Get with non-existing key returns false",
			testfunc: func(t *testing.T) {
				b, _ := setup(t)
				val, ok, err := b.Get([]byte("no exist key"))
				require.NoError(t, err)
				require.False(t, ok)
				require.Nil(t, val)
			},
		},
		{
			name: "Successfully close without removing and recover from the db file",
			testfunc: func(t *testing.T) {
				b, file := setup(t)
				k, v := []byte("key"), []byte("value")
				err := b.Set(k, v)
				require.NoError(t, err)

				err = b.Close(false)
				require.NoError(t, err)

				// recover from the file
				b, err = bbolt.New(file, "", mode)
				require.NoError(t, err)

				res, ok, err := b.Get(k)
				require.NoError(t, err)
				require.True(t, ok)
				require.Equal(t, v, res)
			},
		},
		{
			name: "Successfully close with removing",
			testfunc: func(t *testing.T) {
				b, file := setup(t)
				k, v := []byte("key"), []byte("value")
				err := b.Set(k, v)
				require.NoError(t, err)

				// set remove flag to true
				err = b.Close(true)
				require.NoError(t, err)

				require.NoFileExists(t, file)
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			test.testfunc(tt)
		})
	}
}

func Test_bbolt_AsyncSet(t *testing.T) {
	t.Parallel()

	tempdir := t.TempDir()
	tmpfile := filepath.Join(tempdir, "test.db")
	b, err := bbolt.New(tmpfile, "", mode)
	require.NoError(t, err)

	kv := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4",
		"key5": "val5",
	}

	eg, _ := errgroup.New(context.Background())
	for k, v := range kv {
		b.AsyncSet(eg, []byte(k), []byte(v))
	}

	// wait until all set is done
	err = eg.Wait()
	require.NoError(t, err)

	for k := range kv {
		_, ok, err := b.Get([]byte(k))
		require.NoError(t, err)
		require.True(t, ok)
	}

	err = b.Close(true)
	require.NoError(t, err)
}

// NOT IMPLEMENTED BELOW