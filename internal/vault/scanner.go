package vault

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MarkdownFile struct {
	Name    string
	Path    string
	RelPath string
	Dir     string
	Size    int64
	ModTime time.Time
}

func ScanMarkdown(vaultPath string) ([]MarkdownFile, error) {
	var files []MarkdownFile

	entries, err := os.ReadDir(vaultPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(vaultPath, entry.Name())

		info, err := os.Lstat(entryPath)
		if err != nil {
			continue
		}

		if info.Mode()&os.ModeSymlink != 0 {
			resolved, err := filepath.EvalSymlinks(entryPath)
			if err != nil {
				continue
			}
			found, err := walkForMarkdown(resolved, vaultPath, entry.Name())
			if err != nil {
				continue
			}
			files = append(files, found...)
		} else if info.IsDir() {
			found, err := walkForMarkdown(entryPath, vaultPath, entry.Name())
			if err != nil {
				continue
			}
			files = append(files, found...)
		} else if isMarkdown(entry.Name()) {
			fi, err := os.Stat(entryPath)
			if err != nil {
				continue
			}
			relPath, _ := filepath.Rel(vaultPath, entryPath)
			files = append(files, MarkdownFile{
				Name:    stripExtension(entry.Name()),
				Path:    entryPath,
				RelPath: relPath,
				Dir:     ".",
				Size:    fi.Size(),
				ModTime: fi.ModTime(),
			})
		}
	}

	return files, nil
}

func walkForMarkdown(root, vaultPath, linkName string) ([]MarkdownFile, error) {
	var files []MarkdownFile

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if !isMarkdown(d.Name()) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		relFromRoot, _ := filepath.Rel(root, path)
		relPath := filepath.Join(linkName, relFromRoot)
		dir := filepath.Dir(relPath)

		files = append(files, MarkdownFile{
			Name:    stripExtension(d.Name()),
			Path:    path,
			RelPath: relPath,
			Dir:     dir,
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})

		return nil
	})

	return files, err
}

func isMarkdown(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".md" || ext == ".markdown"
}

func stripExtension(name string) string {
	ext := filepath.Ext(name)
	return strings.TrimSuffix(name, ext)
}
