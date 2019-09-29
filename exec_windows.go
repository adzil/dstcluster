// Copyright (c) 2019 Fadhli Dzil Ikram. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

package main

import (
	"errors"
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

func initCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

func interruptCmd(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return errors.New("cannot interrupt nil process")
	}
	return windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, uint32(cmd.Process.Pid))
}
