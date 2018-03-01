// Copyright 2016 Palantir Technologies, Inc.
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

package integration_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel/framework/artifactresolver"
	"github.com/palantir/godel/framework/pluginapitester"
	"github.com/palantir/godel/pkg/products"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	okgoPluginLocator  = "com.palantir.godel:okgo-plugin:0.1.0"
	okgoPluginResolver = "/Users/nmiyake/.m2/repository/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"

	godelYML = `exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`
)

func TestCompilesInProjectDir(t *testing.T) {
	assetPath, err := products.Bin("compiles-asset")
	require.NoError(t, err)

	projectDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()
	projectDir, err = filepath.EvalSymlinks(projectDir)
	require.NoError(t, err)

	const checkYML = ``

	err = os.MkdirAll(path.Join(projectDir, "godel", "config"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(projectDir, "godel", "config", "godel.yml"), []byte(godelYML), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(projectDir, "godel", "config", "check.yml"), []byte(checkYML), 0644)
	require.NoError(t, err)

	specs := []gofiles.GoFileSpec{
		{
			RelPath: "foo.go",
			Src:     "package foo; foo",
		},
	}

	_, err = gofiles.Write(projectDir, specs)
	require.NoError(t, err)

	outputBuf := &bytes.Buffer{}
	pluginCfg := artifactresolver.LocatorWithResolverConfig{
		Locator: artifactresolver.LocatorConfig{
			ID: okgoPluginLocator,
		},
		Resolver: okgoPluginResolver,
	}
	pluginsParam, err := pluginCfg.ToParam()
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(projectDir)
	require.NoError(t, err)
	defer func() {
		err = os.Chdir(wd)
		require.NoError(t, err)
	}()

	runPluginCleanup, err := pluginapitester.RunAsset(pluginsParam, []string{assetPath}, "check", []string{
		"compiles",
	}, projectDir, false, outputBuf)
	defer runPluginCleanup()
	require.EqualError(t, err, "")

	want := `Running compiles...
foo.go:1:14: expected declaration, found 'IDENT' foo
Finished compiles
`
	assert.Equal(t, want, outputBuf.String())
}

func TestCompilesInInnerProjectDir(t *testing.T) {
	assetPath, err := products.Bin("compiles-asset")
	require.NoError(t, err)

	projectDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()
	projectDir, err = filepath.EvalSymlinks(projectDir)
	require.NoError(t, err)

	const checkYML = ``

	err = os.MkdirAll(path.Join(projectDir, "godel", "config"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(projectDir, "godel", "config", "godel.yml"), []byte(godelYML), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(projectDir, "godel", "config", "check.yml"), []byte(checkYML), 0644)
	require.NoError(t, err)

	specs := []gofiles.GoFileSpec{
		{
			RelPath: "foo.go",
			Src:     "package foo; foo",
		},
	}

	_, err = gofiles.Write(projectDir, specs)
	require.NoError(t, err)

	outputBuf := &bytes.Buffer{}
	pluginCfg := artifactresolver.LocatorWithResolverConfig{
		Locator: artifactresolver.LocatorConfig{
			ID: okgoPluginLocator,
		},
		Resolver: okgoPluginResolver,
	}
	pluginsParam, err := pluginCfg.ToParam()
	require.NoError(t, err)

	innerDir := path.Join(projectDir, "inner")
	err = os.MkdirAll(innerDir, 0755)
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(innerDir)
	require.NoError(t, err)
	defer func() {
		err = os.Chdir(wd)
		require.NoError(t, err)
	}()

	runPluginCleanup, err := pluginapitester.RunAsset(pluginsParam, []string{assetPath}, "check", []string{
		"compiles",
	}, projectDir, false, outputBuf)
	defer runPluginCleanup()
	require.EqualError(t, err, "")

	want := `Running compiles...
../foo.go:1:14: expected declaration, found 'IDENT' foo
Finished compiles
`
	assert.Equal(t, want, outputBuf.String())
}
