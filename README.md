# vmx
Remote instances management tool over SSH written in Go

```
$ vmx
NAME:
   vmx -
vmx is a tool for interacting with cloud instances (like AWS EC2, for example) over SSH


USAGE:
   vmx [global options] command [command options] [arguments...]

VERSION:
   0.0.0

AUTHOR:
   Aliaksandr Kazlou

COMMANDS:
     run      Run custom command
     list     List available custom commands
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Config

All the control files `commands`, `defaults`, `hosts`, etc. are looking for in the `$VMX_HOME` directory if defined,
otherwise in directory `$HOME/.vmx` (so for unix like systems it would be enough to put them in the `~/.vmx`).

## Commands

There are 2 available commands: `list` and `run`.

### `run` command

`run` command is the core of the tool, it is main purpose.

`run` runs the commands described in the `commands` file

#### "commands" file

`commands` is the `ini` file with the following syntax:

```
[command-name]
workingdir=
command=
```

where `command-name` is the command name which will be used in the `run` command, i.e. `vmx run host-name command-name`.

`workingdir` is optional, i.e. you can provide the working dir to change before running the command.

`command` is the required field.

Ex.:

```
[mem]
command=df -h
```

or

```
[rest-logs]
workingdir=/opt/app
command=tail -f -n 10 logs/rest.log
```

which could be run as following:

- `vmx run host-name mem` or
- `vmx run host-name rest-logs`

For the `rest-logs` run command it will change the working directory to the `/opt/apt` first before running the command.

### Command with the confirmation

If you put the `!` in the end of the command name, i.e. `[redeploy!]` then when running this command (by using the
command name without the `!`), the app will ask the confirmation before running the command on the host, i.e.

```
# commands file
...
[redeploy!]
command=./redeploy.sh
...
```

and

```
vmx run host-name redeploy
```

and

```
vmx run dev redeploy
Confirm to run "redeploy" command on [host-name] - yes/no or y/n:
```

## Credits
- [ssh_config](https://github.com/kevinburke/ssh_config)
- [go-ini](https://github.com/go-ini/ini)
- [cli](https://github.com/urfave/cli)
