FROM golang:1.17.1

WORKDIR /app

COPY .. .

EXPOSE 8080-8081

RUN make build

RUN make gen-ca

CMD ./server.out