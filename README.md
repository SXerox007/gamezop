# gamezop

## Run:

'make app' for run the gRPC server

```
make app
```

'make rest' for api exposed by the gRPC server for the REST point

```
make rest
```

## Ports
* gRPC server on localhost:50051
* REST server on localhost:5051
* Mongodb on localhost:27017
* redis on localhost:6379
* Kafka on localhost:9092


## CURL:

```
curl -X POST -k http://localhost:5051/v1/game/gamezop -d '{"player_name":"test","email":"test@gmail.com"}'
```