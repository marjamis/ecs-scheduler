default:

clean:
	rm $(GOPATH)/bin/ecs-scheduler

local_build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(GOPATH)/bin/ecs-scheduler github.com/marjamis/ecs-scheduler

test: # Test all packages within the repo with verbose output
	go test -v -cover $(shell go list ./... | grep -v /vendor/)

mintest: # Test all packages within the repo
	go test -cover ./...
