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

### Build & Run
```shell
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOCACHE=/tmp/.cache/ go build -ldflags="-s -w" -o build/reference-api-mock-server && \
build/reference-api-mock-server ./docker/api/openapi.yaml
```

## Swagger UI

Links:
* [Swagger UI](https://swagger.io/tools/swagger-ui/)
* [Swagger UI GitHub](https://github.com/swagger-api/swagger-ui)

### Start
```shell
docker compose run --rm --service-ports swagger-ui
```

[Local Swagger UI](http://127.0.0.1:8081/)

## RapiDoc

Links:
* [RapiDoc](https://rapidocweb.com/)  
* [RapiDoc Examples](https://rapidocweb.com/list.html)
* [RapiDoc GitHub](https://github.com/rapi-doc/RapiDoc)

### Start
```shell
docker compose run --rm --service-ports rapidoc
```

[Local RapiDoc](http://127.0.0.1:8082/)
