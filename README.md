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
   1.0.0

AUTHOR:
   Aliaksandr Kazlou

COMMANDS:
     run      Run custom command
     list     List available custom commands
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --profile value, -p value  profile to use to read hosts and commands from
   --help, -h                 show help
   --version, -v              print the version
```

Table of Contents
=================

* [Config](#config)
* [Profiles] (#profiles)
* [Commands](#commands)
   * [run command](#run-command)
      * ["commands" file](#commands-file)
      * [Command with the confirmation](#command-with-the-confirmation)
      * [Pass extra arguments to the command](#pass-extra-arguments-to-the-command)
      * [Running the ad-hoc command](#running-the-ad-hoc-command)
   * [list command](#list-command)
* [Hosts](#hosts)
   * [Ad-hoc host name](#ad-hoc-host-name)
* [Defaults](#defaults)
* [Bash auto completion](#bash-auto-completion)
* [Credits](#credits)


## Config

All the control files `commands`, `defaults`, `hosts`, etc. are looking for in the `$VMX_HOME` directory if defined,
otherwise in directory `$HOME/.vmx` (so for unix like systems it would be enough to put them in the `~/.vmx`).

## Profiles

`vmx` supports profiles, similar to the `aws` CLI from Amazon. So you can pass either `-p` or `--profile`
to specify which profile to use. If not specified, it will also check for the value of the env var
`$VMX_DEFAULT_PROFILE`.

If any of those set that profile will be used, otherwise the "default" one.

What profile really means, is that the configuration files will be read from `$VMX_HOME/<profile>`, instead of
`$VMX_HOME`.

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

#### Command with the confirmation

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

#### Pass extra arguments to the command

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

#### Running the ad-hoc command

It is also possible to run the ad-hoc command, i.e. the command which is not defined in the `commands` file.

Ex.:

```
vmx run host-name df -h
```

with no `df` command definition in the `commands` file, will be interpreted as the "ad-hoc" command, and will be
executed on the host as it is, i.e. `df -h`.

### `list` command

`list` command simply prints all available `run` commands, i.e.:

```
$ vmx list

app-logs
check-version
less-logs
logs
mem
redeploy (!)
rest-logs
stop (!)
view-docker-compose
```

## Hosts

Hosts configuration is based on the notion of host groups, exactly the same configuration concept as Ansible has.

Hosts are defined in the `$VMX_HOME/hosts` (or otherwise in `~/.vmx/hosts`), and they are tightly coupled with the
`$VMX_SSH_CONFIG_HOME/config` (or `~/.ssh/config` otherwise).

Here the syntax for the `hosts` file:

```
[group-name]
host1
host2
etc.

[group-name:children]
group-name1
group-name2
```

There is also the special hosts group named `all`.

The actual hostname, user, identity file, etc. are necessary for SSH to know are defined in the
`$VMX_SSH_CONFIG_HOME/config` file.

Example:

```
# ~/.ssh/config
Host rest-prod1
    User ubuntu
    Hostname 1.2.3.4

Host rest-prod2
    User ubuntu
    Hostname 5.6.7.8
    IdentityFile ~/.ssh/rest_prod2_id_rsa
```

and the corresponding section of the `hosts` file:

```
# ~/.vmx/hosts
[rest-prod]
rest-prod1
rest-prod2
```

And so the to execute the `run` command for both `rest-prod` instances:

```
vmx run rest-prod redeploy
```

### Ad-hoc host name

If host name used in the `run` command is not defined in the `hosts` file it is then looked in the `~/.ssh/config`
directly instead. So, if you don't have hosts to group you don't need then to configure `hosts` at all, and utilize yours
`~/.ssh/config` instead.

## Defaults

You can configure default values per host groups in the `$VMX_HOME/defaults`, and also supports `all` hosts group.

The only supported value in the `defaults` is `workingdir`, i.e.:

```
[all]
workingdir=/opt/app

[rest-prod]
workingdir=/opt/app/rest
```

## Bash auto completion

To enable Bash completion, copy `autocomplete/vmx` file into `/etc/bash_completion.d/` directory. Then either restart
the current shell, or `source` that added auto completion file.

## Credits
- [ssh_config](https://github.com/kevinburke/ssh_config)
- [go-ini](https://github.com/go-ini/ini)
- [cli](https://github.com/urfave/cli)
