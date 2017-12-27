package state

import (
	"errors"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

var (
	clusterName = "testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func generateCIARNs(count int, m *mockECSClient) {
	cis := make([]*string, count)
	for i := 0; i < count; i++ {
		cis[i] = aws.String("arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-" + strconv.Itoa(i))
	}

	m.ciARNs = append(m.ciARNs, cis...)
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
	}, m.lciError
}

func (m *mockECSClient) DescribeContainerInstances(*ecs.DescribeContainerInstancesInput) (*ecs.DescribeContainerInstancesOutput, error) {
	return &ecs.DescribeContainerInstancesOutput{}, m.dciError
}

func TestGetInstanceARNs(t *testing.T) {
	t.Run("0 Container Instances", func(t *testing.T) {
		m := &mockECSClient{}
		_, err := getInstanceARNs(&clusterName, m)
		if err != nil {
			assert.Equal(t, errors.New("Function: getInstanceARNs: No Container Instances in Cluster"), err)
		}
	})

	t.Run("2 Container Instances", func(t *testing.T) {
		m := &mockECSClient{}
		generateCIARNs(2, m)
		output, err := getInstanceARNs(&clusterName, m)
		if err != nil {
			t.FailNow()
		}

		for i := 0; i < len(output); i++ {
			log.WithFields(log.Fields{
				"in-memory": *m.ciARNs[i],
				"response":  *output[i],
			}).Debug()
			assert.Equal(t, *m.ciARNs[i], *output[i])
		}
	})

	t.Run("Uses NextToken", func(t *testing.T) {
		m := &mockECSClient{}
		generateCIARNs(200, m)
		output, err := getInstanceARNs(&clusterName, m)
		if err != nil {
			t.FailNow()
		}

		for i := 0; i < len(output); i++ {
			log.WithFields(log.Fields{
				"in-memory": *m.ciARNs[i],
				"response":  *output[i],
			}).Debug()
			assert.Equal(t, *m.ciARNs[i], *output[i])
		}
	})

	t.Run("Error in response", func(t *testing.T) {
		m := &mockECSClient{}
		m.lciError = errors.New("Unknown error")
		_, err := getInstanceARNs(&clusterName, m)
		if err != nil {
			assert.Equal(t, m.lciError, err)
		}
	})

	t.Run("1..1000 Container Instances", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}
		for i := 1; i <= 1000; i++ {
			m := &mockECSClient{}
			generateCIARNs(i, m)

			output, err := getInstanceARNs(&clusterName, m)
			if err != nil {
				t.FailNow()
			}

			for j := 0; j < len(output); j++ {
				log.WithFields(log.Fields{
					"in-memory": *m.ciARNs[j],
					"response":  *output[j],
				}).Debug()
				assert.Equal(t, *m.ciARNs[j], *output[j])
			}
		}
	})
}

func TestDescribeContainerInstances(t *testing.T) {
	t.Run("Error in response from getInstanceARNs", func(t *testing.T) {
		m := &mockECSClient{}
		// Requires at least 1 Container Instance otherwise gets caught in the 0 Container Instance error output
		generateCIARNs(1, m)
		m.dciError = errors.New("Unknown error")
		_, err := DescribeContainerInstances(&clusterName, m)
		if err != nil {
			assert.Equal(t, m.dciError, err)
		}
	})

	t.Run("0 Container Instances", func(t *testing.T) {
		m := &mockECSClient{}
		generateCIARNs(0, m)
		_, err := DescribeContainerInstances(&clusterName, m)
		if err != nil {
			assert.Equal(t, errors.New("Function: getInstanceARNs: No Container Instances in Cluster"), err)
		}
	})

	t.Run("1 Container Instances", func(t *testing.T) {
		m := &mockECSClient{}
		generateCIARNs(1, m)
		_, err := DescribeContainerInstances(&clusterName, m)
		if err != nil {
			t.FailNow()
		}
	})
}
