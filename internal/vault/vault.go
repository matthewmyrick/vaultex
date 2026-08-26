package vault

import (
	"fmt"
	"os"
	"path/filepath"
)

func Create(vaultsDir, name string) error {
	vaultPath := filepath.Join(vaultsDir, name)

	if _, err := os.Stat(vaultPath); err == nil {
		return fmt.Errorf("vault directory %q already exists", vaultPath)
	}

	if err := os.MkdirAll(vaultPath, 0755); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	return nil
}

func Delete(vaultsDir, name string) error {
	vaultPath := filepath.Join(vaultsDir, name)

	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return nil
	}

	if err := os.RemoveAll(vaultPath); err != nil {
		return fmt.Errorf("failed to remove vault directory: %w", err)
	}

	return nil
}

func GetPath(vaultsDir, name string) string {
	return filepath.Join(vaultsDir, name)
}

func Exists(vaultsDir, name string) bool {
	_, err := os.Stat(filepath.Join(vaultsDir, name))
	return err == nil
}
