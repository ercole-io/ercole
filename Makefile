BINARY=ercole
VERSION=`git log --format=%H -n 1`
BUILD=`date -u +%Y-%m-%d-%H:%M:%S-UTC`
LDFLAGS=-ldflags "-X github.com/ercole-io/ercole/v2/cmd.serverVersion=${VERSION}"

run: build
	./ercole -c $(ERCOLE_CONF) serve
conf: build
	./ercole -c $(ERCOLE_CONF) show-config

build:
	CGO_ENABLED=0 go build ${LDFLAGS} -o ${BINARY}

test:
	go clean -testcache
	go test ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	find . -name "fake_*_test.go" -exec rm "{}" \;
	go generate ./...
	go clean -testcache