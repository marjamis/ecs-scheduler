package schedulers

import (
	// "reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/stretchr/testify/assert"
	// log "github.com/sirupsen/logrus"
)

func TestLeastTasks(t *testing.T) {
	t.Run("nil Input", func(t *testing.T) {
		data := LeastTasks(nil)
		assert.Nil(t, data)
	})

	t.Run("0 Instances Input", func(t *testing.T) {
		instances := &ecs.DescribeContainerInstancesOutput{
			ContainerInstances: nil,
		}
		data := LeastTasks(instances)
		assert.Nil(t, data)
	})

	t.Run("1 Container Instances", func(t *testing.T) {
		cis := make([]*ecs.ContainerInstance, 1)
		val := int64(1)
		cis[0] = &ecs.ContainerInstance{
			ContainerInstanceArn: aws.String("arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-1"),
			RunningTasksCount:    &val,
			PendingTasksCount:    &val,
		}
		instances := &ecs.DescribeContainerInstancesOutput{
			ContainerInstances: cis,
		}
		data := LeastTasks(instances)
		assert.Equal(t, "arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-1", *data)
	})

	t.Run("2 Container Instances", func(t *testing.T) {
		cis := make([]*ecs.ContainerInstance, 2)
		val := int64(2)
		val2 := int64(1)
		cis[0] = &ecs.ContainerInstance{
			ContainerInstanceArn: aws.String("arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-1"),
			RunningTasksCount:    &val,
			PendingTasksCount:    &val,
		}
		cis[1] = &ecs.ContainerInstance{
			ContainerInstanceArn: aws.String("arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-2"),
			RunningTasksCount:    &val2,
			PendingTasksCount:    &val2,
		}
		instances := &ecs.DescribeContainerInstancesOutput{
			ContainerInstances: cis,
		}
		data := LeastTasks(instances)
		assert.Equal(t, "arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-2", *data)
	})
}
