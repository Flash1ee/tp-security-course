IMAGE = proxy
CONTAINER = server
FILENAME = server.out
PORT = 8080

.PHONY: build docker-build docker-run docker-stop docker-rm clean

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

gen-ca:
	sh ./cert/gen_ca.sh

rm-certs:
	sudo rm -rf ./certs/

stop:
	docker-compose stop

run:
	docker-compose up --build --no-deps

go-to-container:
	docker exec -it https-proxy_main_1 bash

setup:
	sh ./cert/setup.sh