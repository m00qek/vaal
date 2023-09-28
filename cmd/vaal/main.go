package main

import (
	"errors"

	"github.com/docopt/docopt-go"
)

const (
	usage = `
Prints a shortened version of a directory absolute path.

Usage:
  vaal [--host <host>] [--config <yaml>]... (show | copy) FILE... [--dry-run]
  vaal [--host <host>] [--config <yaml>]... (sync | shell) [--dry-run]
  vaal [--host <host>] [--config <yaml>]... run <command>... [--dry-run]
  vaal --help

Options:
  -d --dry-run        Host name.
  -c --config=<yaml>  Host name.
  -H --host=<host>    Host name.
  -h --help           Show this message.
`
)

func pickHost(config Config, hostname interface{}) (*ConfigHost, error) {
	if hostname == nil {
		if config.DefaultHost == nil {
			return nil, errors.New("There's no default-host configured and none were passed as an argument.")
		}

		hostname = *config.DefaultHost
	}

	host := config.Hosts[hostname.(string)]
	if host == nil {
		return nil, errors.New("Could not find the config entry for host '" + hostname.(string) + "'.")
	}

	return host, nil
}

func main() {
	arguments, _ := docopt.ParseDoc(usage)

	config, err := LoadConfig(arguments["--config"].([]string))
	if err != nil {
		panic(err)
	}

	host, err := pickHost(*config, arguments["--host"])
	if err != nil {
		panic(err)
	}

	dryRun, _ := arguments.Bool("--dry-run")

	if arguments["run"] == true {
		err = Run(*host, dryRun, arguments["<command>"].([]string))
	} else if arguments["shell"] == true {
		err = Shell(*host, dryRun)
	} else if arguments["show"] == true {
		err = Show(*host, arguments["FILE"].([]string))
	} else if arguments["copy"] == true {
		err = Copy(*host, dryRun, arguments["FILE"].([]string))
	} else if arguments["sync"] == true {
		err = Sync(*host, dryRun)
	}

	if err != nil {
		panic(err)
	}
}
