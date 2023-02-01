# cw
cisco worker - small tool that helps automate some work

```
cw - 'Cisco worker' helps automate some stuff

Usage:
  cw [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  do          Executes  commands from the script file on devices
  getconfig   Gets running configs and saves to files
  help        Help about any command
  shusers     Shows users local account on devices

Flags:
      --config string   config file (default is $HOME/.cw.yaml)
  -h, --help            help for cw
      --hosts string    hosts file (default is ./hosts)

Use "cw [command] --help" for more information about a command.
```
TODO:
- add check command to check availability and get hostname
- after connecting, read the prompt, determine the mode($/#) and hostname from promt.
- invent a device model and make it possible to serialize to JSON


DONE:
- set timeout for ssh session


