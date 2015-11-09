# ulule-tools
A set of Go tools to manage [Ulule](http://www.ulule.com/) projects, orders, supporters...

So far, it's a simple a client API package to consume [Ulule's API](http://developers.ulule.com).

To get started real quick, use the Dockerfile and build the app in a container:

```shell
docker build -t ulule-tools .
```

Run the container:

```shell
docker run -ti --rm --name ulule-tools ulule-tools
```
`--rm` is just an optional flag to remove the container when you exit

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







