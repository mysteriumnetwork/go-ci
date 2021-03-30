// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package commands

import (
	"github.com/magefile/mage/sh"
)

// Test runs the tests.
func Test(path string) error {
	return sh.RunV("go", "test", "-race", "-cover", path)
}
