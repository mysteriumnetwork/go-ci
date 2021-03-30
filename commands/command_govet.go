// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package commands

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// GoVet checks that the source is compliant with go vet
func GoVet(path string, additionalArgs ...string) error {
	args := []string{"vet", path}
	args = append(args, additionalArgs...)
	out, err := sh.Output("go", args...)
	fmt.Print(out)
	if err != nil {
		fmt.Println("❌ GoVet")
		return err
	}
	fmt.Println("✅ GoVet")
	return nil
}
