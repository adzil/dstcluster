# Don't Starve Together Cluster Runner

`dstcluster` provides cross-platform command-line interface for running
multiple shards in a Don't Starve Together cluster configuration.

## Quick Start

Before using `dstcluster` you need to install DST dedicated game server using
`steamcmd`, assuming that it is already installed on `/usr/local/steam` and the
game will be installed to `/usr/local/dst`.

```
/usr/local/steam/steamcmd.sh \
    +@ShutdownOnFailedCommand 1 \
    +login anonymous \
    +force_install_dir /usr/local/dst \
    +app_update 343050 validate \
    +quit
```

Then, download `dstcluster` from the [Github release page](https://github.com/adzil/dstcluster/releases)
to the game `bin/` directory. On MacOS, it is located in `dontstarve_dedicated_server_nullrenderer.app/Contents/MacOS/`
instead of `bin/`. The following command can be used to download `dstcluster`
from terminal if you are using Linux.

```
curl -L https://github.com/adzil/dstcluster/releases/download/v1.0.2/dstcluster-linux-amd64.tar.gz | \
    tar xzC /usr/local/dst/bin
```

Now `dstcluster` can be executed the same way as you run the server command.
It will pick up all shard inside a cluster configuration and run it under a
single `dstcluster` instance. For example:

```
/usr/local/dst/bin/dstcluster -cluster MyDediServer
```

## Console Input

Because of there are multiple shard instances running inside `dstcluster`, we
need to be able to switch between shard's `STDIN` to enter a lua command. To do
that simply type `:<shard_name>` followed by a new line. Subsequent `STDIN`\
input after the command will be passed through to the appointed shard. Use the
command again to switch between `STDIN`s.

## Graceful Termination

`dstcluster` can be safely terminated with `SIGINT` (or `Ctrl+C`) or `SIGTERM`
if you are running it under Docker container. On Windows, sending `Ctrl+C` or
`Ctrl+Break` signal will also trigger graceful termination. Under the hood, it
sends `SIGINT` (or `Ctrl+Break` signal on Windows) to every running shards so
it can cleanly terminate itself and wait until the last running shard has been
terminated.

Note that the default Docker termination wait time (10 seconds) may be not
enough for the shard to clean themself up. It is possible to corrupt the save
data if the termination wait time has already passed and Docker forcefully kill
the container.

## Building from Source

To build `dstcluster` from the source you need Go 1.13 installed. Then, simply
clone this repository and run `go build` or `go install` to automatically put
them in `$GOPATH/bin`:

```
git clone https://github.com/adzil/dstcluster.git
cd dstcluster
go build .
```

Because of the game platform compatiblity, `dstcluster` only supports Windows,
MacOS, and Linux.

`make` commands only used to create and compress cross-platform binary for the
[release page](https://github.com/adzil/dstcluster/releases).

## Command-Line Options

This information are directly taken from
[the Klei forum](https://forums.kleientertainment.com/forums/topic/64743-dedicated-server-command-line-options-guide/).
Any shard related command-line options are disabled. Please use `server.ini` to
configure the shard instead.

### -server_path

This option is only applicable to `dstcluster`.

Change the dedicated game server binary path. If it points to a relative path,
`dstcluster` will check them relative to the current working directory and the
executable directory. This option is useful if `dstcluster` does not reside
alongside the game binary. The default value depends on the platform:

Platform | Directory
:------: | ---------
Windows  | `.\dontstarve_dedicated_server_nullrenderer.exe`
MacOS    | `./dontstarve_dedicated_server_nullrenderer`
Linux    | `./dontstarve_dedicated_server_nullrenderer`

If you are using the default value on MacOS, `dstcluster` must be stored in
`dontstarve_dedicated_server_nullrenderer.app/Contents/MacOS` instead of
alongside the `.app` folder.

### -sync_mods

This option is only applicable to `dstcluster`.

If this option is enabled, `dstcluster` will try to populate all mods inside
the cluster's `modoverrides.lua` and append them to the server's
`dedicated_server_mods_setup.lua` so it can be automatically downloaded and/or
updated when starting the cluster.

### -persistent_storage_root

Change the directory that your configuration directory resides in. This must be
an absolute path. The full path to your files will be
`<persistent_storage_root>/<conf_dir>/` where `<conf_dir>` is the value set by
`-conf_dir`. The default for this option depends on the platform:

Platform | Directory
:------: | ---------
Windows  | `<Your documents folder>/Klei`
MacOS    | `<Your home folder>/Documents/Klei`
Linux    | `~/.klei`

### -conf_dir

Change the name of your configuration directory. This name should not contain
any slashes. The full path to your files will be
`<persistent_storage_root>/<conf_dir>` where `<persistent_storage_root>` is the
value set by the `-persistent_storage_root` option. The default is
`DoNotStarveTogether`.

### -cluster

Set the name of the cluster directory that this server will use. The server
will expect to find the cluster.ini file in the following location:
`<persistent_storage_root>/<conf_dir>/<cluster>/cluster.ini`, where
`<persistent_storage_root>` and `<conf_dir>` are the values set by the
`-persistent_storage_root` and `-conf_dir options`. The default is `Cluster_1`.

### -offline

Start the server in offline mode. In offline mode, the server will not be
listed publicly, only players on the local network will be able to join, and
any steam-related functionality will not work.

### -disabledatacollection

Disable data collection for the server.

We require the collection of user data to provide online services. Servers
with disabled data collection will only have access to play in offline mode.
For more details on our privacy policy and how we use the data we collect,
please see our official privacy policy. <https://klei.com/privacy-policy>

### -bind_ip

Change the address that the server binds to when listening for player
connections. This is an advanced feature that most people will not need to use.

### -players

Valid values: `1..64`

Set the maximum number of players that will be allowed to join the game. This
option overrides the `[GAMEPLAY] / max_players` setting in `cluster.ini`.

### -backup_logs

Create a backup of the previous log files each time the server is run. The
backups will be stored in a directory called `backup` in the same directory as
`server.ini`.

### -tick

Valid values: `15..60`

This is the number of times per-second that the server sends updates to
clients. Increasing this may improve precision, but will result in more network
traffic. This option overrides the `[NETWORK] / tick_rate` setting in
`cluster.ini`. It is recommended to leave this at the default value of `15`. If
you do change this option, it is recommended that you do so only for LAN games,
and use a number evenly divisible into `60` (`15`, `20`, `30`).
