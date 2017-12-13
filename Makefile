default:

clean:
	rm ./bin/ecs-scheduler.go

local_build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(GOPATH)bin/ecs-scheduler.go ecs-scheduler

run:
	go run scheduler.go constants.go main.go
