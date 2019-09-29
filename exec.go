// Copyright (c) 2019 Fadhli Dzil Ikram. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

package main

import (
	"os/exec"
)

// Command creates new execute command with OS-specific initialization.
func Command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	initCmd(cmd)
	return cmd
}

// Interrupt sends OS-specific interrupt signal.
func Interrupt(cmd *exec.Cmd) error {
	return interruptCmd(cmd)
}
