 BINARY="prometheus-dummy"
 VERSION=1.0.0
 BUILD=$(date + %FT%T%z)

 PACKAGES=$(go list ./... | grep -v /vendor.)
 VETPACKAGES=$(go list ./... | grep -v /vendor/ | grep -v /examples/)
 GOFILES=$(shell find . -name "*.go" -type f -not -path "./vendor/*")

 default: vet fmt fmt-check test
	@go build -o $(BINARY) -tags=jsoniter

list: 
	@echo ${PACKAGES}
	@echo ${VETPACKAGES}
	@echo ${GOFILES}

fmt:
	@gofmt -s -w $(GOFILES)

fmt-check:
	@ape=$(gofmt -s -d $(GOFILES)); \
	if [ -n "$$ape" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${ape}"; \
		exit 1; \
	fi;

test:
	@go test -cpu=1,2,4 -v -tags integration ./...

vet:
	@go vet $(VETPACKAGES)

docker:
    @docker build -t wuxiaoxiaoshen/example:latest .

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: default fmt fmt-check test vet clean