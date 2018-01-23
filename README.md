# ECS Scheduler

| Type | Product | Badge |
| ---- | ------- | ----- |
| Golang checker | Go Report Card | [![Go Report Card](http://goreportcard.com/badge/marjamis/ecs-scheduler)](http://goreportcard.com/report/marjamis/ecs-scheduler) |
| CI/CD - Build | Semaphore CI | [![Build Status](https://semaphoreci.com/api/v1/marjamis/ecs-scheduler/branches/master/badge.svg)](https://semaphoreci.com/marjamis/ecs-scheduler) |
| CI/CD - Build | Travis CI | [![Build Status](https://travis-ci.org/marjamis/ecs-scheduler.svg?branch=master)](https://travis-ci.org/marjamis/ecs-scheduler) |
| CI/CD - Build | CircleCI |  [![CircleCI](https://circleci.com/gh/marjamis/ecs-scheduler/tree/master.svg?style=svg)](https://circleci.com/gh/marjamis/ecs-scheduler/tree/master) |
| CI/CD - Build | Codeship | [![Coverage Status](https://coveralls.io/repos/github/marjamis/ecs-scheduler/badge.svg?branch=master)](https://coveralls.io/github/marjamis/ecs-scheduler?branch=master) |
| Coverage | Coveralls | [![Coverage Status](https://coveralls.io/repos/github/marjamis/ecs-scheduler/badge.svg?branch=master)](https://coveralls.io/github/marjamis/go_ecs-scheduler?branch=master) |
| Coverge | CodeCov | [![codecov](https://codecov.io/gh/marjamis/ecs-scheduler/branch/master/graph/badge.svg)](https://codecov.io/gh/marjamis/ecs-scheduler) |
| Coverage/Linting | Codacy | [![Codacy Badge](https://api.codacy.com/project/badge/Grade/16a51951592a4671aeb01707f74ad59f)](https://www.codacy.com/app/marjamis/ecs-scheduler?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=marjamis/ecs-scheduler&amp;utm_campaign=Badge_Grade) |
| Project Tracking | Waffle | [![Waffle.io - Columns and their card count](https://badge.waffle.io/marjamis/ecs-scheduler.svg?columns=all)](https://waffle.io/marjamis/ecs-scheduler) |

## To build and run:
1. Clone repository.
2. Get the required dependencies with:
```bash
  dep ensure
```
3. Build the application:
```bash
  make local_build
```
4. Run the application with the required settings:
```bash
  $GOPATH/bin/ecs-scheduler --task-definition <task_definition> --region <region> --cluster <cluster_name>
```
