default:

clean:
	rm $(GOPATH)/bin/ecs-scheduler.go

local_build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(GOPATH)/bin/ecs-scheduler ecs-scheduler

test: # Test all packages within the repo with verbose output
	go test -v -cover ./...

mintest: # Test all packages within the repo
	go test -cover ./...
