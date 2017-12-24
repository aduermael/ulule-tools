# ulule-tools

A set of Dockerized applications for Ulule project managers. Making good use of [aduermael/ulule-api-client](https://github.com/aduermael/ulule-api-client)... ☺️

### Before you run any of the following commands

```shell
export APIKEY=<YOUR API KEY>
export USERNAME=<YOUR USERNAME>
export PROJECTID=<YOUR PROJECT ID>
```

### Lotery

Select one or several winners among contributors. A different amount of tickets can be attributed for each level of contribution.

```shell
# build
docker build -t lotery ./lotery

# run
docker run -ti lotery $APIKEY $USERNAME $PROJECTID

# then it's all interactive! :)
```