# npkill

Find and delete node_modules in target dir.

## Install

```sh
go install github.com/nanlei2000/npkill@latest
```

## Example output

```sh
➜  npkill l
path                                                                 file_count    size
/Users/lielienan/Project/esbuild/require/yarnpnp/bar/node_modules    2             20 B
```

## Manual

```sh
➜  npkill -h
NAME:
   npkill - delete node modules

USAGE:
   npkill [global options] command [command options] [arguments...]

COMMANDS:
   list, l  list all node modules folders
   del, d   delete node modules folder
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```
