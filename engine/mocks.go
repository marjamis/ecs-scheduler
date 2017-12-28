package engine

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

//Define a mock struct to be used in your unit tests
type mockECSClient struct {
	ecsiface.ECSAPI

	ciARNs          []*string
	startTaskOutput ecs.StartTaskOutput

	lciError error
	lciCount int32

	dciError error
	DCIO     *ecs.DescribeContainerInstancesOutput

	stError error
}