package action

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/stretchr/testify/assert"
)

func (m *mockECSClient) StartTask(*ecs.StartTaskInput) (*ecs.StartTaskOutput, error) {
	return &m.startTaskOutput, m.stError
}

func TestStartTask(t *testing.T) {
	m := &mockECSClient{}
	clusterName := "testing"

	t.Run("Error", func(t *testing.T) {
		m.stError = errors.New("Unknown Error")
		instanceARN := "arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-1"
		err := StartTask(&instanceARN, &clusterName, m, "arn:aws:ecs:us-west-2:109951093165:task-definition/website:1")
		assert.Error(t, err)
		m.stError = nil
	})

	failures := make([]*ecs.Failure, 1)
	instanceARN := "arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-1"
	failure := "Because"
	failures[0] = &ecs.Failure{
		Arn:    &instanceARN,
		Reason: &failure,
	}
	m.startTaskOutput.Failures = failures
	t.Run("Returning a failure", func(t *testing.T) {
		instanceARN := "arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-1"
		StartTask(&instanceARN, &clusterName, m, "arn:aws:ecs:us-west-2:109951093165:task-definition/website:1")
	})
	m.startTaskOutput.Failures = nil

	t.Run("Normal functioning", func(t *testing.T) {
		instanceARN := "arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-1"
		err := StartTask(&instanceARN, &clusterName, m, "arn:aws:ecs:us-west-2:109951093165:task-definition/website:1")
		assert.NoError(t, err)
	})
}
