## jdb

```
gofmt -w . && go run --race ./cmd -port 3500

gofmt -w . && go run --race ./cmd -port 3600 -rpc 4600

gofmt -w . && go build --race -a -v -o ./cmd ./jdb
```
