GOLANG_VERSION=1.14

run: build
	./ercole -c $(ERCOLE_CONF) serve
conf: build
	./ercole -c $(ERCOLE_CONF) show-config

build:
	go build -o ercole ./main.go

test:
	go generate ./... && go clean -testcache && go test ./...