// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package env

import (
	"os"

	"github.com/pkg/errors"
)

// EnsureEnvVars ensures that specified environment variables have a value.
// Purpose is to reduce boilerplate in mage target when reading multiple env vars.
func EnsureEnvVars(vars ...BuildVar) error {
	missingVars := make([]BuildVar, 0)
	for _, v := range vars {
		if os.Getenv(string(v)) == "" {
			missingVars = append(missingVars, v)
		}
	}
	if len(missingVars) > 0 {
		return errors.Errorf("the following environment variables are missing: %v", missingVars)
	}
	return nil
}
