# ulule-tools
A set of Go tools to manage [Ulule](http://www.ulule.com/) projects, orders, supporters... using [Ulule's API](http://developers.ulule.com).

To get started real quick, use the Dockerfile to build a Docker image:

```shell
docker build -t ulule-tools .
```

Run the Docker container:

```shell
docker run -ti --rm ulule-tools
```
`--rm` is just an optional flag to remove the container when you exit.

### Ulule cli

A command line interface to list your projects, select one, then list supporters and orders.

Build the CLI and run it:

```shell
# from /go
cd src/cmd/cli
go build
./cli
```

You'll be prompted for a username and API key. You can also launch the cli with these arguments:

```shell
./cli <username> <apikey>
```
Commands:

```shell
> project list
> project select <id or slug>
> project supporters
> project orders
```
### Ulule sync

A tool to store orders associated with one project in a local redis database. It's then super fast to get stats from that local snapshot using **statorders** tool (see below).

```shell
# from /go
cd src/cmd/sync
go build
./sync
# enter username
username>
# enter api key (from settings > privacy)
apikey>
# projects will be listed here
# enter project id to sync information
project id>
# enter a name for following sync operation
sync name>
# be patient while syncing rewards & orders
```




