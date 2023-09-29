package main

import (
	"fmt"
)

func ListHosts(config Config) error {
	for key, _ := range config.Hosts {
		fmt.Println(key)
	}

	return nil
}

func Get(host ConfigHost, expr string) error {
	rendered, err := EvalExpression(host, expr)
	if err != nil {
		return err
	}

	if *rendered != "<no value>" {
		fmt.Println(*rendered)
	}

	return nil
}

func Show(host ConfigHost, paths []string) error {
	for _, path := range paths {
		rendered, err := ReadAndApplyTemplate(host, path)
		if err != nil {
			return err
		}

		fmt.Println(*rendered)
	}

	return nil
}

func Shell(host ConfigHost, dryRun bool) error {
	return SecureShell(host, dryRun)
}

func Run(host ConfigHost, dryRun bool, command []string) error {
	return SecureShell(host, dryRun, command...)
}

func Copy(host ConfigHost, dryRun bool, paths []string) error {
	return WithTempDir(dryRun, func(tempdir string) error {
		for _, path := range paths {
			rendered, err := ReadAndApplyTemplate(host, path)
			if err != nil {
				return err
			}

			absPath, relPath, err := CreateTempFile(host, tempdir, path, *rendered)
			SecureCopyToRemote(host, dryRun, absPath, relPath)
		}
		return nil
	})
}

func Sync(host ConfigHost, dryRun bool) error {
	paths, err := ListAllValidFiles(host)
	if err != nil {
		return err
	}

	return Copy(host, dryRun, paths)
}
