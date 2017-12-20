package engine

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
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

	t.Run("Normal values nothing configured", func(t *testing.T) {
		debug = true
		leastTasksSched = true
		cluster = "testing"
		region = "us-west-2"
		taskDefinition = "testing:1"
		assert.Equal(t, ExitStateError, Run())
	})

	// t.Run("Normal values nothing configured", func(t *testing.T) {
	// 	debug = true
	// 	leastTasksSched = true
	// 	cluster = "testing"
	// 	region = "us-west-2"
	// 	taskDefinition = "testing:1"
	// 	assert.Equal(t, ExitNoValidContainerInstance, Run())
	// })
}
