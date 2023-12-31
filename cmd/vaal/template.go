package main

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

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

	contents := "" // header
	for _, path := range bundle {
		text, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		contents += string(text) + "\n"
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

func EvalExpression(config ConfigHost, expr string) (*string, error) {
	return applyTemplate(config, expr, "{{ "+expr+" }}")
}

func ReadAndApplyTemplate(config ConfigHost, templatePath string) (*string, error) {
	templateText, err := readTemplate(config, templatePath)
	if err != nil {
		return nil, err
	}

	rendered, err := applyTemplate(config, templatePath, *templateText)
	if err != nil {
		return nil, err
	}

	return rendered, nil
}
