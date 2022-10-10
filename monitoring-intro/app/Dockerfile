FROM golang:1.19

WORKDIR /go/src/app
COPY . .

RUN go mod tidy
RUN go build -o app

CMD ["/go/src/app/app"]