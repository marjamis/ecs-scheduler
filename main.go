package main

import (
	"ecs-scheduler/engine"

	"os"
)

func main() {
	os.Exit(engine.Run())
}
