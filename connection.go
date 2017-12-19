package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func connectToECS(region string) (svc *ecs.ECS) {
	return ecs.New(session.New(), aws.NewConfig().WithRegion(region))
}
