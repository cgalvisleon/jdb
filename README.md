## JDB

## Sytem vars

- NODEID:

- DB_NAME:

- DB_DRIVER:

- DB_HOST:

- DB_PORT:

- DB_USER:

- DB_PASSWORD:

- APP_NAME:

## Dependencis

```
go get github.com/cgalvisleon/et/@v0.1.3
```

## Run and build

```
gofmt -w . && go run --race ./cmd/jdb -port 3500

gofmt -w . && go run --race ./cmd/jdb -port 3600 -rpc 4600
gofmt -w . && go run --race ./cmd -port 3600 -rpc 4600

gofmt -w . && go build --race -a -o ./jdb ./cmd/jdb
gofmt -w . && go build --race -a -v -o ./jdb ./cmd/jdb

ps aux | grep jdb | grep -v grep

```

## Version

```
Version: v0.1.3
```

## Go work

```
go work init ../et
```

## Kube

```
k -n prd describe pod apigateway
k -n prd get deployments
k -n prd rollout history deployment apigateway
k -n prd get configmap
kubectl describe node
```
