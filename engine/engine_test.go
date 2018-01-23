package engine

import (
	"errors"
	"testing"

	"github.com/marjamis/ecs-scheduler/mocks"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func TestConnectToECS(t *testing.T) {
	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()

	for _, p := range partitions {
		for id := range p.Regions() {
			assert.Equal(t, id, *ConnectToECS(&id).Client.Config.Region)
		}
	}
}

func TestRun(t *testing.T) {
	t.Run("No input values", func(t *testing.T) {
		assert.Equal(t, ExitInvalidCLIOptions, Run())
	})

	t.Run("All input values", func(t *testing.T) {
		debug = true
		scheduler = "leastTasks"
		cluster = "testing"
		region = "us-west-2"
		taskDefinition = "testing:1"
		assert.Equal(t, ExitStateError, Run())
	})

	t.Run("Normal values nothing configured", func(t *testing.T) {
		m := &mocks.MockECSClient{}
		debug = true
		scheduler = "leastTasks"
		cluster = "testing"
		region = "us-west-2"
		taskDefinition = "testing:1"
		assert.Equal(t, ExitStateError, selectProcess(m))
	})

	t.Run("Normal values nothing configured but startTask error", func(t *testing.T) {
		cis := make([]*ecs.ContainerInstance, 1)
		val := int64(1)
		cis[0] = &ecs.ContainerInstance{
			ContainerInstanceArn: aws.String("arn:aws:ecs:us-west-2:101234567891:container-instance/11111111-11a7-469d-b903-1"),
			RunningTasksCount:    &val,
			PendingTasksCount:    &val,
		}
		instances := &ecs.DescribeContainerInstancesOutput{
			ContainerInstances: cis,
		}
		m := &mocks.MockECSClient{
			DCIO: instances,
		}
		m.StError = errors.New("Here I am")
		mocks.GenerateCiARNs(2, m)
		debug = true
		scheduler = "leastTasks"
		cluster = "testing"
		region = "us-west-2"
		taskDefinition = "testing:1"
		assert.Equal(t, ExitStartTaskFailure, selectProcess(m))
	})

	t.Run("Normal values nothing configured", func(t *testing.T) {
		cis := make([]*ecs.ContainerInstance, 1)
		val := int64(1)
		cis[0] = &ecs.ContainerInstance{
			ContainerInstanceArn: aws.String("arn:aws:ecs:us-west-2:101234567891:container-instance/11111111-11a7-469d-b903-1"),
			RunningTasksCount:    &val,
			PendingTasksCount:    &val,
		}
		instances := &ecs.DescribeContainerInstancesOutput{
			ContainerInstances: cis,
		}
		m := &mocks.MockECSClient{
			DCIO: instances,
		}
		mocks.GenerateCiARNs(2, m)
		debug = true
		scheduler = "leastTasks"
		cluster = "testing"
		region = "us-west-2"
		taskDefinition = "testing:1"
		assert.Equal(t, ExitSuccess, selectProcess(m))
	})

	t.Run("Normal values nothing configured", func(t *testing.T) {
		instances := &ecs.DescribeContainerInstancesOutput{
			ContainerInstances: nil,
		}
		m := &mocks.MockECSClient{
			DCIO: instances,
		}
		mocks.GenerateCiARNs(2, m)
		debug = true
		scheduler = "leastTasks"
		cluster = "testing"
		region = "us-west-2"
		taskDefinition = "testing:1"
		assert.Equal(t, ExitNoValidContainerInstance, selectProcess(m))
	})
}
