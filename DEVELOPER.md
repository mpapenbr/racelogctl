# How to release

## Preparations
Install the prerequisites  
```shell
./devtools.sh
```

## Check
This will perform a local goreleaser build in order to check if goreleaser is happy with provided config.

```shell
goreleaser --snapshot --skip-publish --rm-dist
```

## Create a release
A new release is done by calling the script
```shell
./newRelease.sh
```

You may pass arguments to this script. These will forwarded to the tool [svu](https://github.com/caarlos0/svu). In case you want automatic semantic versioning by the latests commits just omit any parameters. Otherwise see the svu documentation.




# Get event data from server

## leipert 3h sebring race

get event json
```shell
go run main.go --url wss://crossbar.leipert-esports.de/ws event info 2  --format json --pretty > leipert-event-2-info.json
```

go run main.go --url wss://crossbar.leipert-esports.de/ws event states 2 --from 1647025565 --full --num 100 --output leipert-event-2-data.txt

## leipert 12h sebring race

get event json
```shell
go run main.go --url wss://crossbar.leipert-esports.de/ws event info 3  --format json --pretty > leipert-event-3-info.json
```

go run main.go --url wss://crossbar.leipert-esports.de/ws event states 3 --from 1648298696 --full --num 200 --output leipert-event-3-data.txt