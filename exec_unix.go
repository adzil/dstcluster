// Copyright (c) 2019 Fadhli Dzil Ikram. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

// +build !windows

package main

import (
	"errors"
	"os/exec"
	"syscall"
)

func initCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func interruptCmd(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return errors.New("cannot interrupt nil process")
	}
	return cmd.Process.Signal(syscall.SIGINT)
}
