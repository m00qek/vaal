package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const configDir = "etc/config"

func relativePath(source string, file string) (string, error) {
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	return filepath.Rel(source, path)
}

func WithTempDir(dryRun bool, fn func(dirpath string) error) error {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}

	if !dryRun {
		defer os.RemoveAll(tempDir)
	}

	return fn(tempDir)
}

func CreateTempFile(host ConfigHost, tempdir string, path string, content string) (string, string, error) {
	relpath, err := relativePath(*host.Source, path)
	if err != nil {
		return "", "", err
	}

	abspath := filepath.Join(tempdir, relpath)

	os.MkdirAll(filepath.Dir(abspath), os.ModePerm)

	tempFile, err := os.Create(abspath)
	if err != nil {
		return "", "", err
	}
	defer tempFile.Close()

	_, err = tempFile.WriteString(content)
	if err != nil {
		return "", "", err
	}

	return abspath, relpath, nil
}

func ListAllValidFiles(host ConfigHost) ([]string, error) {
	paths := []string{}

	err := filepath.WalkDir(*host.Source, func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relpath, err := relativePath(*host.Source, path)
		if !strings.HasPrefix(relpath, configDir) {
			file, err := os.Stat(path)
			if err != nil {
				return err
			}

			if !file.IsDir() {
				paths = append(paths, path)
			}
		} else if filepath.Dir(relpath) == configDir {
			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}
