# ECS Scheduler
[![Build Status](https://travis-ci.org/marjamis/ecs-scheduler.svg?branch=master)](https://travis-ci.org/marjamis/ecs-scheduler)
[![Coverage Status](https://coveralls.io/repos/github/marjamis/ecs-scheduler/badge.svg?branch=master)](https://coveralls.io/github/marjamis/ecs-scheduler?branch=master)
[![CircleCI](https://circleci.com/gh/marjamis/ecs-scheduler.svg?style=svg)](https://circleci.com/gh/marjamis/ecs-scheduler)
[![Codeship Status for marjamis/ecs-scheduler](https://app.codeship.com/projects/42c1f6a0-ee70-0134-6d5c-62b847b8d86d/status?branch=master)](https://app.codeship.com/projects/208712)

## To build and run:
1. Clone repository.
1. Get required additional libraries with:

  make setup

1. Run the application. Basic version such as:

  go run src/ecs-scheduler/scheduler.go --task-definition \<task_definition\> --region \<region\> --cluster \<cluster_name\>

   and if you're happy with the current state compile a statically compiled version with:

  make build
