// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package job

import (
	"log"
	"os"
)

// Precondition skips the job (exits successfully) if the precondition fails.
func Precondition(p func() bool) {
	if !p() {
		log.Println("Precondition failed, skipping the job")
		os.Exit(0)
	}
	log.Println("Precondition passed")
}
