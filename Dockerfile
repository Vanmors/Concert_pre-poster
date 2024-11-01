FROM golang:latest

RUN go version
ENV GOPATH=/
COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres. sh executable
RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o concert-pre-poster ./cmd/main.go

CMD ["./concert-pre-poster"]

EXPOSE 8000
