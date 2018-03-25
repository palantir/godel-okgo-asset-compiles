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

package main

import (
	"os"

	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/okgo/checker"
	"github.com/palantir/pkg/cobracli"

	"github.com/palantir/godel-okgo-asset-compiles/compiles/config"
	"github.com/palantir/godel-okgo-asset-compiles/compiles/creator"
	"github.com/palantir/godel-okgo-asset-compiles/generated_src"
)

func main() {
	os.Exit(amalgomated.RunApp(os.Args, nil, amalgomated.NewCmdLibrary(amalgomatedcheck.Instance()), checkMain))
}

func checkMain(osArgs []string) int {
	os.Args = osArgs
	var debugFlagVal bool
	rootCmd := checker.AssetRootCmd(creator.Compiles(), config.UpgradeConfig, "run compiles check")
	return cobracli.ExecuteWithDefaultParamsWithVersion(rootCmd, &debugFlagVal, "")
}
