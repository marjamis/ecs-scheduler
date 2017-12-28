package engine

import (
	"errors"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func (m *mockECSClient) StartTask(*ecs.StartTaskInput) (*ecs.StartTaskOutput, error) {
	return &m.startTaskOutput, m.stError
}

func (m *mockECSClient) ListContainerInstances(input *ecs.ListContainerInstancesInput) (*ecs.ListContainerInstancesOutput, error) {
	var cis []*string
	var token *string

	if *input.MaxResults <= 100 && int64(len(m.ciARNs)) <= 100 {
		cis = m.ciARNs
	} else if input.NextToken != nil {
		tokenint, err := strconv.Atoi(*input.NextToken)
		if err != nil {
			return nil, err
		}
		tokenNum := int64(tokenint)
		var total int64

		if (tokenNum + *input.MaxResults) > int64(len(m.ciARNs)) {
			total = int64(len(m.ciARNs))
		} else {
			total = tokenNum + *input.MaxResults
		}

		cis = m.ciARNs[tokenNum:total]

		if (tokenNum + *input.MaxResults) < int64(len(m.ciARNs)) {
			val := tokenNum + *input.MaxResults
			string2 := strconv.FormatInt(int64(val), 10)
			token = &string2
		}
	} else if *input.MaxResults < int64(len(m.ciARNs)) {
		cis = m.ciARNs[:*input.MaxResults]
		string2 := strconv.FormatInt(*input.MaxResults, 10)
		token = &string2
	}

	return &ecs.ListContainerInstancesOutput{
		ContainerInstanceArns: cis,
		NextToken:             token,
	}, nil
}

func (m *mockECSClient) DescribeContainerInstances(*ecs.DescribeContainerInstancesInput) (*ecs.DescribeContainerInstancesOutput, error) {
	return m.DCIO, nil
}

func generateCIARNs(count int, m *mockECSClient) {
	cis := make([]*string, count)
	for i := 0; i < count; i++ {
		cis[i] = aws.String("arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-" + strconv.Itoa(i))
	}

	m.ciARNs = append(m.ciARNs, cis...)
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

	//Need to fix this shitty test
	t.Run("All input values", func(t *testing.T) {
		debug = true
		scheduler = "leastTasks"
		cluster = "testing"
		region = "us-west-2"
		taskDefinition = "testing:1"
		assert.Equal(t, ExitStateError, Run())
	})

	t.Run("Normal values nothing configured", func(t *testing.T) {
		m := &mockECSClient{}
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
			ContainerInstanceArn: aws.String("arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-1"),
			RunningTasksCount:    &val,
			PendingTasksCount:    &val,
		}
		instances := &ecs.DescribeContainerInstancesOutput{
			ContainerInstances: cis,
		}
		m := &mockECSClient{
			DCIO: instances,
		}
		m.stError = errors.New("Here I am")
		generateCIARNs(2, m)
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
			ContainerInstanceArn: aws.String("arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-1"),
			RunningTasksCount:    &val,
			PendingTasksCount:    &val,
		}
		instances := &ecs.DescribeContainerInstancesOutput{
			ContainerInstances: cis,
		}
		m := &mockECSClient{
			DCIO: instances,
		}
		generateCIARNs(2, m)
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
		m := &mockECSClient{
			DCIO: instances,
		}
		generateCIARNs(2, m)
		debug = true
		scheduler = "leastTasks"
		cluster = "testing"
		region = "us-west-2"
		taskDefinition = "testing:1"
		assert.Equal(t, ExitNoValidContainerInstance, selectProcess(m))
	})
}
