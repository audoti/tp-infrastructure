Topics:

 - aeroport/wind

## Install

Download dependencies

```
go get ./...
```

## Run

Start fakeiot

```
go run cmd/fakeiot/main.go
```

Start Mosquitto broker (port **18830**)

```
./scripts/startMosquittoBroker.sh
# or
docker run -it -p 18830:1883 -p 9001:9001 eclipse-mosquitto
```

Start Redis

```
./scripts/startRedis.sh
# or
docker run -it -p 6379:6379 redis
# or
redis-server
```
