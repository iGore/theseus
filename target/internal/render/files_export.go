package render

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func WriteOutputFile(outPath string, content string) error {
	dir := filepath.Dir(outPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := os.WriteFile(outPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	return nil
}

func ExportLicenseFiles(modules ModuleResultMap, targetDir string) error {
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("create export directory: %w", err)
	}
	for _, key := range sortedKeys(modules) {
		m := modules[key]
		if m.LicenseFile == "" {
			continue
		}
		if err := copyFile(m.LicenseFile, filepath.Join(targetDir, key, filepath.Base(m.LicenseFile))); err != nil {
			return fmt.Errorf("export license file for %s: %w", key, err)
		}
	}
	return nil
}

func copyFile(srcPath string, dstPath string) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
