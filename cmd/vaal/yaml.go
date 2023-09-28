package main

import (
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

type YAML = map[string]interface{}

func getIn[T any](dict YAML, keys ...string) *T {
	var current_dict interface{}
	current_dict = dict

	for _, key := range keys[:len(keys)-1] {
		v := current_dict.(YAML)
		current_dict = v[key]

		if current_dict == nil {
			return nil
		}
	}

	v := current_dict.(YAML)
	key := keys[len(keys)-1]

	if v[key] == nil {
		return nil
	}

	value := v[key].(T)
	return &value
}

func updateIn(value interface{}, dict *YAML, keys ...string) {
	var current_dict interface{}
	current_dict = *dict

	for _, key := range keys[:len(keys)-1] {
		v := current_dict.(YAML)
		current_dict = v[key]

		if current_dict == nil {
			return
		}
	}

	v := current_dict.(YAML)
	key := keys[len(keys)-1]
	v[key] = value
}

func loadYaml(file string) (YAML, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config YAML
	err = yaml.Unmarshal(content, &config)

	if err != nil {
		return nil, err
	}

	if config["hosts"] != nil {
		hosts := config["hosts"].(YAML)

		for _, host := range hosts {
			bHost := host.(YAML)
			source := getIn[string](bHost, "source")
			if source != nil {
				absSource, err := filepath.Abs(filepath.Join(filepath.Dir(file), *source))
				if err != nil {
					return nil, err
				}

				updateIn(absSource, &bHost, "source")
			}
		}
	}

	return config, nil
}

func fuseYaml(config1 YAML, config2 YAML) YAML {
	fused := YAML{}

	for k, v1 := range config1 {
		v2 := config2[k]

		if v1 == nil {
			fused[k] = v2
			continue
		}

		if v2 == nil {
			fused[k] = v1
			continue
		}

		if v1m, ok := v1.(YAML); ok {
			if v2m, ok := v2.(YAML); ok {
				fused[k] = fuseYaml(v1m, v2m)
				continue
			}
		}

		fused[k] = v2
	}

	for k, v2 := range config2 {
		if fused[k] == nil {
			fused[k] = v2
		}
	}

	return fused
}

func LoadYAMLs(files []string) (YAML, error) {
	configs := []YAML{}

	for _, file := range files {
		config, err := loadYaml(file)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	fused := configs[0]
	for _, config := range configs[1:] {
		fused = fuseYaml(fused, config)
	}

	return fused, nil
}
