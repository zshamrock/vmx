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

### Pass extra arguments to the command

Any extra arguments passed to the the `run` command will be then passed to the actual command to run. Here the example:

```
# commands file
...
[logs]
workingdir=/opt/app/logs
command=tail -n 10
...
```

And when running:

```
vmx run host-name logs -f rest.log
```

it will be interpreted as `tail -n 10 -f rest.log` (i.e. all extra arguments are passed to the `command` defined in the
`commands` file).

### Running the ad-hoc command

It is also possible to run the ad-hoc command, i.e. the command which is not defined in the `commands` file.

Ex.:

```
vmx run host-name df -h
```

with no `df` command definition in the `commands` file, will be interpreted as the "ad-hoc" command, and will be
executed on the host as it is, i.e. `df -h`.

## Credits
- [ssh_config](https://github.com/kevinburke/ssh_config)
- [go-ini](https://github.com/go-ini/ini)
- [cli](https://github.com/urfave/cli)
