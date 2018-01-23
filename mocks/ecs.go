package mocks

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

//MockECSClient Define a mock struct to be used in your unit tests. Easily configurable to the unit tests needs.
type MockECSClient struct {
	ecsiface.ECSAPI
	CiARNs []*string

	LciError error
	LciCount int32

	DciError        error
	StartTaskOutput ecs.StartTaskOutput
	StError         error
	DCIO            *ecs.DescribeContainerInstancesOutput
}

//GenerateCiARNs Will generate data for faking ContainerInstance details
func GenerateCiARNs(count int, m *MockECSClient) {
	cis := make([]*string, count)
	for i := 0; i < count; i++ {
		cis[i] = aws.String("arn:aws:ecs:us-west-2:101234567891:container-instance/11111111-11a7-469d-b903-" + strconv.Itoa(i))
	}

	m.CiARNs = append(m.CiARNs, cis...)
}

//ListContainerInstances A form of the ECS ListContainerInstances API call
func (m *MockECSClient) ListContainerInstances(input *ecs.ListContainerInstancesInput) (*ecs.ListContainerInstancesOutput, error) {
	var cis []*string
	var token *string

	if *input.MaxResults <= 100 && int64(len(m.CiARNs)) <= 100 {
		cis = m.CiARNs
	} else if input.NextToken != nil {
		tokenint, err := strconv.Atoi(*input.NextToken)
		if err != nil {
			return nil, err
		}
		tokenNum := int64(tokenint)
		var total int64

		if (tokenNum + *input.MaxResults) > int64(len(m.CiARNs)) {
			total = int64(len(m.CiARNs))
		} else {
			total = tokenNum + *input.MaxResults
		}

		cis = m.CiARNs[tokenNum:total]

		if (tokenNum + *input.MaxResults) < int64(len(m.CiARNs)) {
			val := tokenNum + *input.MaxResults
			string2 := strconv.FormatInt(int64(val), 10)
			token = &string2
		}
	} else if *input.MaxResults < int64(len(m.CiARNs)) {
		cis = m.CiARNs[:*input.MaxResults]
		string2 := strconv.FormatInt(*input.MaxResults, 10)
		token = &string2
	}

	return &ecs.ListContainerInstancesOutput{
		ContainerInstanceArns: cis,
		NextToken:             token,
	}, m.LciError
}

//DescribeContainerInstances A form of the ECS DescribeContainerInstances API call
func (m *MockECSClient) DescribeContainerInstances(*ecs.DescribeContainerInstancesInput) (*ecs.DescribeContainerInstancesOutput, error) {
	return m.DCIO, m.DciError
}

//StartTask A form of the ECS StartTask API call
func (m *MockECSClient) StartTask(*ecs.StartTaskInput) (*ecs.StartTaskOutput, error) {
	return &m.StartTaskOutput, m.StError
}
