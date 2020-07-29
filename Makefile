GOLANG_VERSION=1.14

run: build
	./ercole -c $(ERCOLE_CONF) serve
conf: build
	./ercole -c $(ERCOLE_CONF) show-config

build:
	go build -o ercole ./main.go

test:
	go clean -testcache
	go test ./...

clean:
	rm -rf ercole
	find . -name "fake_*_test.go" -exec rm "{}" \;
	go generate ./...
	go clean -testcache