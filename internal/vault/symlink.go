package vault

import (
	"fmt"
	"os"
	"path/filepath"
)

func AddLink(vaultPath, targetPath, linkName string) error {
	info, err := os.Stat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("target path %q does not exist", targetPath)
		}
		return fmt.Errorf("failed to stat target path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("target path %q is not a directory", targetPath)
	}

	symPath := filepath.Join(vaultPath, linkName)
	if _, err := os.Lstat(symPath); err == nil {
		return fmt.Errorf("link %q already exists in vault", linkName)
	}

	if err := os.Symlink(targetPath, symPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

func RemoveLink(vaultPath, linkName string) error {
	symPath := filepath.Join(vaultPath, linkName)

	info, err := os.Lstat(symPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("link %q does not exist in vault", linkName)
		}
		return fmt.Errorf("failed to stat link: %w", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("%q is not a symlink", linkName)
	}

	if err := os.Remove(symPath); err != nil {
		return fmt.Errorf("failed to remove symlink: %w", err)
	}

	return nil
}

type LinkStatus struct {
	Name    string
	Target  string
	Valid   bool
	SymPath string
}

func ListLinks(vaultPath string) ([]LinkStatus, error) {
	entries, err := os.ReadDir(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault directory: %w", err)
	}

	var links []LinkStatus
	for _, entry := range entries {
		entryPath := filepath.Join(vaultPath, entry.Name())
		info, err := os.Lstat(entryPath)
		if err != nil {
			continue
		}

		if info.Mode()&os.ModeSymlink == 0 {
			continue
		}

		target, err := os.Readlink(entryPath)
		if err != nil {
			continue
		}

		_, statErr := os.Stat(entryPath)
		links = append(links, LinkStatus{
			Name:    entry.Name(),
			Target:  target,
			Valid:   statErr == nil,
			SymPath: entryPath,
		})
	}

	return links, nil
}
