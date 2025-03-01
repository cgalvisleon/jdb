## jdb

```
gofmt -w . && go run --race ./cmd/jdb -port 3500

gofmt -w . && go run --race ./cmd/jdb -port 3600 -rpc 4600
gofmt -w . && go run --race ./cmd -port 3600 -rpc 4600

gofmt -w . && go build --race -a -o ./jdb ./cmd/jdb
gofmt -w . && go build --race -a -v -o ./jdb ./cmd/jdb

ps aux | grep jdb | grep -v grep

```

## Library

```
go get github.com/cgalvisleon/et/@v1.0.3
```

## Go work

```
go work init ../et
```
