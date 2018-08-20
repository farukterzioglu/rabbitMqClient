# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=publisherWebApp
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME) -consulUrl=https://demo.consul.io/ui/dc1/kv -consulPath=rabbitMqConsumerGoLang
deps:
	$(GOGET) github.com/farukterzioglu/rabbitMqClient
	$(GOGET) github.com/gorilla/mux

# Cross compilation
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/farukterzioglu/rabbitMqClient/PublisherWebApp golang:latest go build -o "$(BINARY_UNIX)" -v