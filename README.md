# ECS Scheduler
[![Build Status](https://travis-ci.org/marjamis/ecs-scheduler.svg?branch=master)](https://travis-ci.org/marjamis/ecs-scheduler)
[![Coverage Status](https://coveralls.io/repos/github/marjamis/ecs-scheduler/badge.svg?branch=master)](https://coveralls.io/github/marjamis/ecs-scheduler?branch=master)
[![CircleCI](https://circleci.com/gh/marjamis/ecs-scheduler.svg?style=svg)](https://circleci.com/gh/marjamis/ecs-scheduler)
[ ![Codeship Status for marjamis/ecs-scheduler](https://app.codeship.com/projects/ae328790-a0ce-0135-44ea-2622f92aca11/status?branch=master)](https://app.codeship.com/projects/254078)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/16a51951592a4671aeb01707f74ad59f)](https://www.codacy.com/app/marjamis/ecs-scheduler?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=marjamis/ecs-scheduler&amp;utm_campaign=Badge_Grade)

## To build and run:
1. Clone repository.
1. Get required additional libraries with:
<<<<<<< HEAD
```bash
  make setup
```
1. Run the application. Basic version such as:
```bash
  go run src/ecs-scheduler/scheduler.go --task-definition \<task_definition\> --region \<region\> --cluster \<cluster_name\>
```
   and if you're happy with the current state compile a statically compiled version with:
```bash
  make build
```
=======

  make setup

1. Run the application. Basic version such as:

  go run src/ecs-scheduler/scheduler.go --task-definition \<task_definition\> --region \<region\> --cluster \<cluster_name\>

   and if you're happy with the current state compile a statically compiled version with:

  make build

## TODO
* Do one of the above with a go plugin - https://golang.org/pkg/plugin/
* ensure contexts are available thought the application
>>>>>>> master
