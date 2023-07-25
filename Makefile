linter:
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.53.3 golangci-lint run -v

tests:
	go test ./internal/...

benchmark:
	go test -bench=. ./internal/general/challenger/... -run=^#

build: build-server build-client

build-server:
	docker build -t word-of-wisdom-server -f ./build/server.Dockerfile .
build-client:
	docker build -t word-of-wisdom-client -f ./build/client.Dockerfile .

start: start-server start-client

start-server:
	echo "starting word-of-wisdom-server"
	docker run --rm -d --env-file ./build/server.env --network host --name word-of-wisdom-server word-of-wisdom-server:latest
	echo "word-of-wisdom-server started"

start-client:
	echo "starting word-of-wisdom-client"
	docker run --rm --env-file ./build/client.env --network host --name word-of-wisdom-client word-of-wisdom-client:latest

stop-server:
	docker rm -f word-of-wisdom-server