# Reference Api Mock Server

## Mock Server

### Run
```shell
go run main.go ./docker/api/openapi.yaml
```

### Build
```shell
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOCACHE=/tmp/.cache/ go build -ldflags="-s -w" -o build/reference-api-mock-server
```

### Run Build
```shell
build/reference-api-mock-server ./docker/api/openapi.yaml
```




## Swagger UI

### Start
```shell
docker compose up -d
```

### Use
[Swagger UI](http://127.0.0.1:8081/)

### Stop

```shell
docker compose down
```