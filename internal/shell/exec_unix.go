//go:build darwin || linux

package shell

import (
	"path/filepath"
	"syscall"
)

func execShell(shellPath string, env []string) error {
	return syscall.Exec(shellPath, []string{filepath.Base(shellPath), "-i"}, env)
}
