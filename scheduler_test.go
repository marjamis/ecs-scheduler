package main

import (
	"errors"
	// "reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.WarnLevel)
}

// Define a mock struct to be used in your unit tests
type mockECSClient struct {
	ecsiface.ECSAPI
	ciARNs []*string

	lciError error
	lciCount int32

	dciError error
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
			log.Info(val)
			string2 := strconv.FormatInt(int64(val), 10)
			token = &string2
		}
	} else if *input.MaxResults < int64(len(m.ciARNs)) {
		cis = m.ciARNs[:*input.MaxResults]
		string2 := strconv.FormatInt(*input.MaxResults, 10)
		token = &string2
	}

	log.WithFields(log.Fields{
		"len(cis)": len(cis),
		"tokenNum": token,
	}).Debug("here")

	return &ecs.ListContainerInstancesOutput{
		ContainerInstanceArns: cis,
		NextToken:             token,
	}, m.lciError
}

func (m *mockECSClient) DescribeContainerInstances(*ecs.DescribeContainerInstancesInput) (*ecs.DescribeContainerInstancesOutput, error) {
	return &ecs.DescribeContainerInstancesOutput{}, m.dciError
}

func TestGetInstanceARNs(t *testing.T) {
	m := &mockECSClient{}
	// if testing.Short() {
	//     t.Skip("skipping test in short mode.")
	// }

	t.Run("0 Container Instances", func(t *testing.T) {
		_, err := getInstanceARNs("testing", m)
		if err != nil {
			assert.Equal(t, errors.New("Function: getInstanceARNs: No Container Instances in Cluster"), err)
		}
	})

	generateCIARNs(2, m)
	t.Run("2 Container Instances", func(t *testing.T) {
		output, err := getInstanceARNs("testing", m)
		if err != nil {
			t.FailNow()
		}

		for i := 0; i < len(output); i++ {
			assert.Equal(t, *m.ciARNs[i], *output[i])
		}
	})

	generateCIARNs(200, m)
	t.Run("Uses NextToken", func(t *testing.T) {
		output, err := getInstanceARNs("testing", m)
		if err != nil {
			t.FailNow()
		}

		for i := 0; i < len(output); i++ {
			log.WithFields(log.Fields{
				"inmemory": *m.ciARNs[i],
				"response": *output[i],
			}).Debug("here")
			assert.Equal(t, *m.ciARNs[i], *output[i])
		}
	})

	m.lciError = errors.New("Unknown error")
	t.Run("Error in response", func(t *testing.T) {
		_, err := getInstanceARNs("testing", m)
		if err != nil {
			assert.Equal(t, m.lciError, err)
		}
	})
}

func TestDescribeContainerInstances(t *testing.T) {
	m := &mockECSClient{}

	m.lciError = errors.New("Unknown error")
	t.Run("Error in response from getInstanceARNs", func(t *testing.T) {
		_, err := describeContainerInstances("testing", m)
		if err != nil {
			assert.Equal(t, m.lciError, err)
		}
	})

	cis := make([]*string, 2)
	cis[0] = aws.String("arn:aws:ecs:us-west-2:109951093165:container-instance/00000000-11a7-469d-b903-6587a1f11bca")
	cis[1] = aws.String("arn:aws:ecs:us-west-2:109951093165:container-instance/11111111-11a7-469d-b903-6587a1f11bca")
	m.ciARNs = append(m.ciARNs, cis...)
	m.lciError = nil
	m.dciError = errors.New("Unknown error")
	t.Run("Error in response", func(t *testing.T) {
		_, err := describeContainerInstances("testing", m)
		if err != nil {
			assert.Equal(t, m.dciError, err)
		}
	})

	// m.err = nil
	// t.Run("Normal response", func(t *testing.T) {
	// 	output, err := describeContainerInstances("testing", m)
	// 	if err != nil {
	// 		t.FailNow()
	// 	}
	//
	// 	assert.Equal(t, reflect.ValueOf(&ecs.DescribeContainerInstancesOutput{}), reflect.ValueOf(output))
	//
	// })

	// What about when more than 100 CIs to describe
	// t.Run("Error in response", func(t *testing.T) {
	// 	_, err := describeContainerInstances("testing", m)
	// 	if err != nil {
	// 		assert.Equal(t, m.err, err)
	// 	}
	// })
}
