package main

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

const header = `
###############################################################################
##          File generated automatically, please do *not* edit it.           ##
###############################################################################
`

func separator(path string) string {
	return "\n## Source: " + path + "\n"
}

func readTemplate(config ConfigHost, path string) (*string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var bundle []string
	if fileInfo.IsDir() {
		dirfiles, err := os.ReadDir(path)
		if err != nil {
			return nil, nil
		}

		for _, file := range dirfiles {
			if !file.IsDir() {
				bundle = append(bundle, filepath.Join(path, file.Name()))
			}
		}
	} else {
		bundle = []string{path}
	}

	contents := header
	for _, path := range bundle {
		text, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		relpath, err := relativePath(*config.Source, path)
		if err != nil {
			return nil, err
		}

		contents += separator(relpath) + string(text)
	}

	return &contents, nil
}

func applyTemplate(config ConfigHost, name string, tplt string) (*string, error) {
	data := map[string]interface{}{
		"server": map[string]interface{}{
			"addr": config.Server.Addr,
			"user": config.Server.User,
		},
		"params": config.Params,
	}

	t, err := template.New(name).Parse(tplt)
	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return nil, err
	}

	result := tpl.String()
	return &result, nil
}

func ReadAndApplyTemplate(config ConfigHost, templateFile string) (*string, error) {
	templateText, err := readTemplate(config, templateFile)
	if err != nil {
		return nil, err
	}

	rendered, err := applyTemplate(config, templateFile, *templateText)
	if err != nil {
		return nil, err
	}

	return rendered, nil
}
