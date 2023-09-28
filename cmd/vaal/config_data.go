package main

type ConfigServer struct {
	Addr *string
	User *string
}

type ConfigUCI struct {
	ToServices   map[string][]string
	Dependencies map[string][]string
}

type ConfigHost struct {
	Source *string
	Server ConfigServer
	Params YAML
}

type Config struct {
	DefaultHost *string
	UCI         ConfigUCI `yaml:"uci"`
	Hosts       map[string]*ConfigHost
}

func hostConfigFromYAML(yaml YAML) ConfigHost {
	host := ConfigHost{}

	host.Source = getIn[string](yaml, "source")
	host.Server.Addr = getIn[string](yaml, "server", "addr")
	host.Server.User = getIn[string](yaml, "server", "user")

	params := getIn[YAML](yaml, "params")
	if params != nil {
		host.Params = *params
	} else {
		host.Params = YAML{}
	}

	return host
}

func mapOfStringArraysFromYAML(yaml *YAML) map[string][]string {
	result := map[string][]string{}

	if yaml == nil {
		return result
	}

	for key, value := range *yaml {
		list := []string{}
		for _, item := range value.([]interface{}) {
			list = append(list, item.(string))
		}

		result[key] = list
	}
	return result
}

func FromYAML(yaml YAML) Config {
	config := Config{}

	config.DefaultHost = getIn[string](yaml, "default-host")

	config.UCI.ToServices = mapOfStringArraysFromYAML(getIn[YAML](yaml, "uci", "to-services"))
	config.UCI.Dependencies = mapOfStringArraysFromYAML(getIn[YAML](yaml, "uci", "dependencies"))

	config.Hosts = map[string]*ConfigHost{}
	hosts := getIn[YAML](yaml, "hosts")
	if hosts != nil {
		for key, value := range *hosts {
			host := hostConfigFromYAML(value.(YAML))
			config.Hosts[key] = &host
		}
	}

	return config
}
