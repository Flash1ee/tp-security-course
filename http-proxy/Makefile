IMAGE = http_proxy
CONTAINER = server
FILENAME = server.out
PORT = 8080

.PHONY: build docker-build docker-run docker-stop docker-rm clean
lint:
	golangci-lint run -c golangci-lint.yml ./...

build:
	go build -o $(FILENAME) ./cmd/main.go

docker-build:
	docker build --no-cache -f ./docker/Dockerfile . --tag $(IMAGE)

docker-run:
	docker run -p $(PORT):$(PORT) --name $(CONTAINER) -t $(IMAGE)

docker-stop:
	docker stop $(CONTAINER)

docker-rm:
	docker rm -vf $(CONTAINER) || true

clean:
	rm -rf $(FILENAME)