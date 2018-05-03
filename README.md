godel-okgo-asset-compiles
=========================
godel-okgo-asset-compiles is an asset for the gödel [okgo plugin](https://github.com/palantir/okgo). It provides the functionality of the [go-compiles](https://github.com/palantir/go-compiles) check.

This check verifies that all of the provided packages (and their tests) compiles. The output is similar to that of `go build ./...`, but unlike that command it also checks for compilation errors in test files for the packages.
