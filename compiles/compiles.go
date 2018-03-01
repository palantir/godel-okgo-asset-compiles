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

package compiles

import (
	"github.com/palantir/okgo/checker"
	"github.com/palantir/okgo/okgo"
)

const (
	TypeName okgo.CheckerType     = "compiles"
	Priority okgo.CheckerPriority = 0
)

func Creator() checker.Creator {
	return checker.NewCreator(
		TypeName,
		Priority,
		func(cfgYML []byte) (okgo.Checker, error) {
			return checker.NewAmalgomatedChecker(TypeName, checker.Priority(Priority)), nil
		},
	)
}
