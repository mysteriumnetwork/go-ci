# CI tools for golang

This is a repo containing common magefile based CI tools for golang projects.

To use it, include a makefile in root of your repo directory

```
# This Makefile is meant to be used by people that do not usually work with Go source code.
# If you know what GOPATH is then you probably don't need to bother with make.

MAGE_PATH=${GOPATH}/src/github.com/magefile/mage
MAGE=go run ${GOPATH}/src/github.com/mysteriumnetwork/go-ci/mage.go -d ./ci

default:
ifeq ("$(wildcard $(MAGE_PATH))","")
	go get -u -d github.com/magefile/mage
	go get github.com/mysteriumnetwork/go-ci
endif
	${MAGE} -l

% :
ifeq ("$(wildcard $(MAGE_PATH))","")
	go get -u -d github.com/magefile/mage
	go get github.com/mysteriumnetwork/go-ci
endif
	${MAGE} $@
```


Then, create a ci folder to contain all the mage files you need. To include the common scripts from this library, a following file is suggested:


```

// Runs the test suite against the repo
func Test() error {
	return commands.Test("../...")
}

// Checks for copyright headers in files
func CheckCopyright() error {
	return commands.Copyright("../...")
}

// Installs go dependencies
func Dep() error {
	return commands.Deps()
}

// Checks for issues with go imports
func CheckGoImports() error {
	return commands.GoImports("../...")
}

// Reports linting errors in the solution
func CheckGoLint() error {
	return commands.GoLint("../...")
}

// Updates the go report for the repo
func CheckGoReport() error {
	return commands.GoReport("github.com/your-name/your-repo")
}

// Checks that the source is compliant with go vet
func CheckGoVet() error {
	return commands.GoVet("../...")
}

// Checks that the source is compliant with go vet
func Check() error {
	return commands.Check("../...")
}

```

With this, just run make in the root of your repo and you're set!
