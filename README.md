# CI tools for golang

This is a repo containing common magefile based CI tools for golang projects.

To use it, include a makefile in root of your repo directory

```
# This Makefile is meant to be used by people that do not usually work with Go source code.
# If you know what GOPATH is then you probably don't need to bother with make.

MAGE_PATH=${GOPATH}/src/github.com/magefile/mage
MAGE=go run ${GOPATH}/src/github.com/mysteriumnetwork/goci/mage.go -d ./ci

default:
ifeq ("$(wildcard $(MAGE_PATH))","")
	go get -u -d github.com/magefile/mage
	go get github.com/mysteriumnetwork/goci
endif
	${MAGE} -l

% :
ifeq ("$(wildcard $(MAGE_PATH))","")
	go get -u -d github.com/magefile/mage
	go get github.com/mysteriumnetwork/goci
endif
	${MAGE} $@
```


Then, create a ci folder to contain all the mage files you need. To include the common scripts from this library, a following file is suggested:


```
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/mysteriumnetwork/goci/commands"
)

// ExludedDirs contains the list of directories that we'll exclude for the repo
var ExludedDirs = []string{".git", "vendor", "build", "docs"}

// Runs the test suite against the repo
func Test() error {
	return commands.Test()
}

// Checks for copyright headers in files
func Copyright() error {
	return commands.Copyright(ExludedDirs)
}

// Installs go dependencies
func Dep() error {
	return commands.Deps()
}

// Checks for issues with go imports
func GoImports() error {
	return commands.GoImports(ExludedDirs)
}

// Reports linting errors in the solution
func GoLint() error {
	return commands.GoLint(ExludedDirs)
}

// Updates the go report for the repo
func GoReport() error {
	return commands.GoReport("github.com/your-name/your-project")
}

// Checks that the source is compliant with go vet
func GoVet() error {
	return commands.GoVet()
}

```

With this, just run make in the root of your repo and you're set!
