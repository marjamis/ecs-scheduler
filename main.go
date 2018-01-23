package main

import (
	"os"

	"github.com/marjamis/ecs-scheduler/engine"
)

func main() {
	os.Exit(engine.Run())
}
