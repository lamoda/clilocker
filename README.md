CLI Locker
==========

**Work is in heavy progress!**

Console utility that helps to control number of commands running at the same time.

This utility could helpful in the following cases:
1. You have a command that must be run in just a given number of copies
2. You have a command that runs on different servers or docker containers and you want to control number of copies

How it works
------------

Cli Locker is a console utility that takes you command as an argument and controls number of it's launches. 
It uses local (file) or remote (redis) locks to control if limit of allowed runs is not exceeded.

Example:
```shell script
clilocker --id=mycommand --limit=3 ls -la
```

Above command will control that command `ls -la` will run in no more than 3 copy at the same time

Cluster - is all clilocker instances that trying to launch the same CLI command with id `mycommand`

Internally on every start of `clilocker` it takes a lock using local files or redis based on the configuration given.

Locks
-----

CliLocker supports three kind of locks:
1. local ones - based on files
2. redis - based on Redis locks
3. redis-sentinel - same as Redis but with Sentinel configured

Configuration
-------------

CliLocker by-default uses file locks.

You can define it's configuration using file `.clilocker.yaml`

Example is below:
```yaml
locker_dsn: redis://redis_address:6379
```

`locker_dsn` - defines lock used to limit number of the commands running

`command_id` - id of the command

`limit` - maximum limit of commands running at the same time 

Cluster locks must be available for the all nodes.

`.clilocker.yaml` file could be placed in the following locations
1. current directory
2. using `--config` flag

CLI Options
-----------

`--limit` - limits number of command instances allowed to run at the same time. Default: no limit

`--config` - path to the configuration file

`--id` - id of the command. Default: command itself

`--locker_dsn` - defines lock used to limit number of the commands running

How to get CliLocker
--------------------

1. Clone this repository
2. Run
```shell script
make build
```