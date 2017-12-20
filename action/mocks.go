package action

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

//Define a mock struct to be used in your unit tests
type mockECSClient struct {
	ecsiface.ECSAPI
	stError         error
	startTaskOutput ecs.StartTaskOutput
}
