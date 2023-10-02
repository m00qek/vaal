vaal
===

Manage configuration files of OpenWRT hosts.

## Dependencies

Make sure that `ssh` and `scp` commands are available on the client machine.

## Installing

```bash
make install
```

## Usage

**WARNING:** It's still early days and this is under development; things might change.

`vaal` allows you to maintain all your OpenWRT configuration files in a git repo,
organized in a sane way. Let's use the [example](./example) in this repo to,
explore some commands:

```bash
cd example/
vaal show router/etc/config/network
```

It will print the `network` configuration already interpolated with values from
the YAML config files. If there is an OpenWRT installation on the IP defined in
the  YAML config files you may also use `vaal copy` to send single files to the
router or `vaal sync` to send all files at once.

## Configuration 

Every time you run `vaal` it checks if the current directory has a `config.yaml`
and, optionally, a `config.secrets.yaml`. Those files contain the information on 
how to locate the routers and values that will be interpolated on the UCI config
files.
The `config.secrets.yaml` should not be committed to git repos and can hold
sensible information like Wireguard secrets, etc. It's contents take precedence
over `config.yaml`, meaning that if the same key is defined in both `vaal` will
use the value in `config.secrets.yaml`.

You can also explicitly pass multiple config files, like

```bash
vaal --config config.yaml --config $HOME/.vaal.yaml list hosts
```

and the rightmost config will have higher precedence. Be aware that when you pass
config files explicitly the program will not load `config.yaml` and `config.secrets.yaml`
automatically.
