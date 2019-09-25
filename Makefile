BINARY_NAME = "sts"

build:
	env GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o $(BINARY_NAME)

run:
	env GO111MODULE=on go build -mod=vendor -o $(BINARY_NAME)
	./$(BINARY_NAME)

test:
	env GO111MODULE=on go test -mod=vendor -v ./...

test-coverage:
	go get github.com/mattn/goveralls && \
	env GO111MODULE=on go test -mod=vendor -v -cover -coverprofile ~/coverage.out.tmp ./... && \
	cat ~/coverage.out.tmp | grep -v "_mock.go" > ~/coverage.out && \
	/usr/local/go/bin/goveralls -coverprofile ~/coverage.out -service=circle-ci -repotoken=$COVERALLS_TOKEN

clean:
	rm -f $(BINARY_NAME)
	env GO111MODULE=on go clean -mod=vendor