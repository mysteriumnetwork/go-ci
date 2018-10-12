# CI tools for golang

This is a repo containing common magefile based CI tools for golang projects.

To use it, include a makefile in root of your repo directory

```makefile
# This Makefile is meant to be used by people that do not usually work with Go source code.
# If you know what GOPATH is then you probably don't need to bother with make.

GO_PATH=$(shell go env GOPATH)
DEP_PATH=$(GO_PATH)/bin/dep
MAGE=go run ci/mage.go

default:
ifeq ("$(wildcard $(DEP_PATH))", "")
	go get -u github.com/golang/dep/cmd/dep
endif
	${DEP_PATH} ensure
	${MAGE} -l

% :
ifeq ("$(wildcard $(DEP_PATH))", "")
	go get -u github.com/golang/dep/cmd/dep
endif
	${DEP_PATH} ensure
	${MAGE} $(MAKECMDGOALS)

```


Then, create a `magefile.go` file to contain all the mage files you need. To include the common scripts from this library, a following file is suggested:

```golang
// Runs the test suite against the repo
func Test() error {
	return commands.Test("./...")
}

// Checks for copyright headers in files
func CheckCopyright() error {
	return commands.Copyright("./...")
}

// Checks for issues with go imports
func CheckGoImports() error {
	return commands.GoImports("./...")
}

// Reports linting errors in the solution
func CheckGoLint() error {
	return commands.GoLint("./...")
}

// Updates the go report for the repo
func CheckGoReport() error {
	return commands.GoReport("github.com/your-name/your-repo")
}

// Checks that the source is compliant with go vet
func CheckGoVet() error {
	return commands.GoVet("./...")
}

// Checks that the source is compliant with go vet
func Check() error {
	return commands.Check("./...")
}

```

With this, just run make in the root of your repo and you're set!
