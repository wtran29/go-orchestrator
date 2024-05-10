# Build an Orchestator in Go
```
The purpose of this project is to learn the ins and outs of how an orchestrator works.
An orchestrator is a system that provides automation for deploying, scaling,
and managing containers. An orchestration system has the following components 
- task, job, scheduler, manager, worker, cluster, and cli
```
## Task

The task is the smallest unit of work in an orchestration system and typically
runs in a container, like a process that runs on a single machine. A task can be
in one of five states: Pending, Scheduled, Running, Completed, or Failed.

## Job

The job is an aggregation of grouped tasks to perform a set of functions.

## Scheduler

The scheduler decides what machine can best host the tasks defined in the job.

## Manager

The manager is the brain of an orchestrator and the main entry point for
users. Users submits their jobs to the manager and the manager 
uses the scheduler to find a machine where the job's task can run.

## Worker

The worker does the heavy lifting of the orchestrator by running the tasks assigned
by the manager.

## Cluster

The cluster is a the logical grouping of all the above components that could be run
from a single physical or virtual machine.

## CLI

CLI is the command line interface that will allow a user to start/stop tasks,
get status of tasks, see the state of machines, start the manager, start the worker