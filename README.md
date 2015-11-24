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

### Ulule CLI

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
### Ulule syncnames

A tool to list all orders and store associated informations like names & emails in a redis database. A webpage can then be generated, for supporters who want to update these informations. These changes are local to the redis database. It may just be useful in some specific use cases, like when we want to display supporter's names on a product. The webpage shows an opt-out checkbox also, for those who don't want to appear on the product.

An email can be sent to each contributor, with a link to their page. (a SendGrid account is required for that feature)





