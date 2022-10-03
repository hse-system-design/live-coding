FROM golang:1.19

COPY go.mod /server/go.mod
COPY main.go /server/main.go
WORKDIR /server

RUN go mod tidy
RUN apt update
RUN apt install -y curl

CMD ["go", "run", "main.go"]