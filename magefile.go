// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// +build mage

package main

import "github.com/mysteriumnetwork/go-ci/commands"

// Check runs all of the checks.
func Check() error {
	return commands.CheckD("./...")
}
