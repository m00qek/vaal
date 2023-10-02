package main

import (
	"errors"

	"github.com/docopt/docopt-go"
)

const (
	usage = `
Manage config files of OpenWRT hosts.

Usage:
  vaal [--host <host>] [--config <yaml>]... copy FILE... [--to-directory=<dir>] [--dry-run]
  vaal [--host <host>] [--config <yaml>]... sync [--to-directory=<dir>] [--dry-run]
  vaal [--host <host>] [--config <yaml>]... show FILE... [--dry-run]
  vaal [--host <host>] [--config <yaml>]... shell [--dry-run]
  vaal [--host <host>] [--config <yaml>]... run <command>... [--dry-run]
  vaal [--host <host>] [--config <yaml>]... get <expr>
  vaal [--config <yaml>]... list hosts
  vaal --help

Options:
  -d --dry-run             Do not change the host, only print actions.
  -t --to-directory=<dir>  Put files in <dir> instead of the the host.
  -c --config=<yaml>       YAML config file.
  -H --host=<host>         Host name.
  -h --help                Show this message.
`
)

func pickHost(config Config, hostname *string) (*ConfigHost, error) {
	if hostname == nil && config.DefaultHost == nil {
		return nil, errors.New("There's no default-host configured and none were passed as an argument.")
	}

	if hostname == nil {
		hostname = config.DefaultHost
	}

	host := config.Hosts[*hostname]
	if host == nil {
		return nil, errors.New("Could not find the config entry for host '" + *hostname + "'.")
	}

	return host, nil
}

func main() {
	arguments, _ := docopt.ParseDoc(usage)
	config, err := LoadConfig(arguments["--config"].([]string))
	if err != nil {
		panic(err)
	}

	var hostname *string
	if v, err := arguments.String("--host"); err == nil {
		hostname = &v
	}

	var dryRun bool
	if v, err := arguments.Bool("--dry-run"); err == nil {
		dryRun = v
	}

	var toDir *string
	if v, err := arguments.String("--to-directory"); err == nil {
		toDir = &v
	}

	host, err := pickHost(*config, hostname)
	if err != nil {
		panic(err)
	}

	if arguments["run"] == true {
		err = Run(*host, dryRun, arguments["<command>"].([]string))
	} else if arguments["shell"] == true {
		err = Shell(*host, dryRun)
	} else if arguments["show"] == true {
		err = Show(*host, arguments["FILE"].([]string))
	} else if arguments["copy"] == true {
		err = Copy(*host, dryRun, arguments["FILE"].([]string), toDir)
	} else if arguments["sync"] == true {
		err = Sync(*host, dryRun, toDir)
	} else if arguments["get"] == true {
		err = Get(*host, arguments["<expr>"].(string))
	} else if arguments["list"] == true && arguments["hosts"] == true {
		err = ListHosts(*config)
	}

	if err != nil {
		panic(err)
	}
}
