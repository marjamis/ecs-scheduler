# ECS Scheduler
[![Build Status](https://travis-ci.org/marjamis/go_ecs-scheduler.svg?branch=master)](https://travis-ci.org/marjamis/go_ecs-scheduler)
[![Coverage Status](https://coveralls.io/repos/github/marjamis/go_ecs-scheduler/badge.svg?branch=master)](https://coveralls.io/github/marjamis/go_ecs-scheduler?branch=master)
[![CircleCI](https://circleci.com/gh/marjamis/ecs-scheduler/tree/master.svg?style=svg)](https://circleci.com/gh/marjamis/ecs-scheduler/tree/master)
[![Codeship Status for marjamis/go_ecs-scheduler](https://app.codeship.com/projects/42c1f6a0-ee70-0134-6d5c-62b847b8d86d/status?branch=master)](https://app.codeship.com/projects/208712)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/16a51951592a4671aeb01707f74ad59f)](https://www.codacy.com/app/marjamis/ecs-scheduler?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=marjamis/ecs-scheduler&amp;utm_campaign=Badge_Grade)
[![Build Status](https://semaphoreci.com/api/v1/marjamis/ecs-scheduler/branches/master/badge.svg)](https://semaphoreci.com/marjamis/ecs-scheduler)

## To build and run:
1. Clone repository.
2. Get the required dependencies with:

  dep ensure

3. Build the application:

  make local_build

4. Run the application with the required settings:

  $GOPATH/bin/ecs-scheduler --task-definition \<task_definition\> --region \<region\> --cluster \<cluster_name\>
