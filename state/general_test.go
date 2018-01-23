package state

import (
	"errors"
	"testing"

	"github.com/marjamis/ecs-scheduler/mocks"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

var (
	clusterName = "testing"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func TestGetInstanceARNs(t *testing.T) {
	t.Run("0 Container Instances", func(t *testing.T) {
		m := &mocks.MockECSClient{}
		_, err := getInstanceARNs(&clusterName, m)
		if err != nil {
			assert.Equal(t, errors.New("Function: getInstanceARNs: No Container Instances in Cluster"), err)
		}
	})

	t.Run("2 Container Instances", func(t *testing.T) {
		m := &mocks.MockECSClient{}
		mocks.GenerateCiARNs(2, m)
		output, err := getInstanceARNs(&clusterName, m)
		if err != nil {
			t.FailNow()
		}

		for i := 0; i < len(output); i++ {
			log.WithFields(log.Fields{
				"in-memory": *m.CiARNs[i],
				"response":  *output[i],
			}).Debug()
			assert.Equal(t, *m.CiARNs[i], *output[i])
		}
	})

	t.Run("Uses NextToken", func(t *testing.T) {
		m := &mocks.MockECSClient{}
		mocks.GenerateCiARNs(200, m)
		output, err := getInstanceARNs(&clusterName, m)
		if err != nil {
			t.FailNow()
		}

		for i := 0; i < len(output); i++ {
			log.WithFields(log.Fields{
				"in-memory": *m.CiARNs[i],
				"response":  *output[i],
			}).Debug()
			assert.Equal(t, *m.CiARNs[i], *output[i])
		}
	})

	t.Run("Error in response", func(t *testing.T) {
		m := &mocks.MockECSClient{}
		m.LciError = errors.New("Unknown error")
		_, err := getInstanceARNs(&clusterName, m)
		if err != nil {
			assert.Equal(t, m.LciError, err)
		}
	})

	t.Run("1..1000 Container Instances", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}
		for i := 1; i <= 1000; i++ {
			m := &mocks.MockECSClient{}
			mocks.GenerateCiARNs(i, m)

			output, err := getInstanceARNs(&clusterName, m)
			if err != nil {
				t.FailNow()
			}

			for j := 0; j < len(output); j++ {
				log.WithFields(log.Fields{
					"in-memory": *m.CiARNs[j],
					"response":  *output[j],
				}).Debug()
				assert.Equal(t, *m.CiARNs[j], *output[j])
			}
		}
	})
}

func TestDescribeContainerInstances(t *testing.T) {
	t.Run("Error in response from getInstanceARNs", func(t *testing.T) {
		m := &mocks.MockECSClient{}
		// Requires at least 1 Container Instance otherwise gets caught in the 0 Container Instance error output
		mocks.GenerateCiARNs(1, m)
		m.DciError = errors.New("Unknown error")
		_, err := DescribeContainerInstances(&clusterName, m)
		if err != nil {
			assert.Equal(t, m.DciError, err)
		}
	})

	t.Run("0 Container Instances", func(t *testing.T) {
		m := &mocks.MockECSClient{}
		mocks.GenerateCiARNs(0, m)
		_, err := DescribeContainerInstances(&clusterName, m)
		if err != nil {
			assert.Equal(t, errors.New("Function: getInstanceARNs: No Container Instances in Cluster"), err)
		}
	})

	t.Run("1 Container Instances", func(t *testing.T) {
		m := &mocks.MockECSClient{}
		mocks.GenerateCiARNs(1, m)
		_, err := DescribeContainerInstances(&clusterName, m)
		if err != nil {
			t.FailNow()
		}
	})
}
