package main

import (
	"errors"
	"os"
	"path/filepath"
)

func readConfigFiles(paths []string) (*Config, error) {
	for _, file := range paths {
		_, err := os.Stat(file)
		if err != nil {
			return nil, err
		}
	}

	yaml, err := LoadYAMLs(paths)
	if err != nil {
		return nil, err
	}

	config := FromYAML(yaml)
	return &config, nil
}

func readDefaultConfigFiles() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	defaultPaths := []string{
		filepath.Join(cwd, "config.yaml"),
		filepath.Join(cwd, "config.secrets.yaml"),
	}

	paths := []string{}
	for _, path := range defaultPaths {
		_, err := os.Stat(path)
		if err == nil {
			paths = append(paths, path)
		}
	}

	if len(paths) == 0 {
		return nil, errors.New("Could not find default config files.")
	}

	yaml, err := LoadYAMLs(paths)
	if err != nil {
		return nil, err
	}

	config := FromYAML(yaml)
	return &config, nil
}

func LoadConfig(paths []string) (*Config, error) {
	var config *Config
	var err error

	if len(paths) == 0 {
		config, err = readDefaultConfigFiles()
	} else {
		config, err = readConfigFiles(paths)
	}

	return config, err
}
