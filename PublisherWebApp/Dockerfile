FROM golang:latest AS builder
RUN go version

COPY . /go/src/github.com/farukterzioglu/rabbitMqClient/PublisherWebApp/
WORKDIR /go/src/github.com/farukterzioglu/rabbitMqClient/PublisherWebApp/
RUN set -x && \
    go get github.com/golang/dep/cmd/dep && \
    dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/farukterzioglu/rabbitMqClient/PublisherWebApp/app .

EXPOSE 8123
ENTRYPOINT ["./app"]