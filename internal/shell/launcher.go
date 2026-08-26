package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func LaunchShell(vaultName, vaultPath string) error {
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		shellPath = "/bin/zsh"
	}

	if err := os.Chdir(vaultPath); err != nil {
		return fmt.Errorf("failed to change directory to vault: %w", err)
	}

	env := os.Environ()
	env = append(env, "VAULTEX_VAULT="+vaultName)
	env = append(env, "VAULTEX_ROOT="+vaultPath)

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		return execShell(shellPath, env)
	}

	return execShellFallback(shellPath, env)
}

func execShellFallback(shellPath string, env []string) error {
	cmd := exec.Command(shellPath, "-i")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	return cmd.Run()
}

func GetActiveVault() (name string, root string, active bool) {
	name = os.Getenv("VAULTEX_VAULT")
	root = os.Getenv("VAULTEX_ROOT")
	active = name != "" && root != ""
	return
}

func DetectEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}

	for _, name := range []string{"nvim", "vim", "vi"} {
		if path, err := exec.LookPath(name); err == nil {
			return filepath.Base(path)
		}
	}

	return "vi"
}
