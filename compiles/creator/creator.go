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

package creator

import (
	"io"
	"strings"

	"github.com/palantir/okgo/checker"
	"github.com/palantir/okgo/okgo"

	"github.com/palantir/godel-okgo-asset-compiles/compiles"
)

func Compiles() checker.Creator {
	return checker.NewCreator(
		compiles.TypeName,
		compiles.Priority,
		func(cfgYML []byte) (okgo.Checker, error) {
			return &wrappedChecker{
				checker: checker.NewAmalgomatedChecker(compiles.TypeName, checker.ParamPriority(compiles.Priority),
					checker.ParamLineParserWithWd(
						func(line, wd string) okgo.Issue {
							if line == "-: " {
								// skip lines that have empty output
								return okgo.Issue{}
							}
							return okgo.NewIssueFromLine(line, wd)
						},
					),
				),
			}, nil
		},
	)
}

type wrappedChecker struct {
	checker okgo.Checker
}

func (w *wrappedChecker) Type() (okgo.CheckerType, error) {
	return w.checker.Type()
}

func (w *wrappedChecker) Priority() (okgo.CheckerPriority, error) {
	return w.checker.Priority()
}

func (w *wrappedChecker) Check(pkgPaths []string, projectDir string, stdout io.Writer) {
	// trim "./" prefixes to support package path formats for Go modules
	for i := range pkgPaths {
		pkgPaths[i] = strings.TrimPrefix(pkgPaths[i], "./")
	}
	w.checker.Check(pkgPaths, projectDir, stdout)
}

func (w *wrappedChecker) RunCheckCmd(args []string, stdout io.Writer) {
	w.checker.RunCheckCmd(args, stdout)
}
