package cron

import (
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/pkg/errors"
	"time"
)

type TaskType string
type runFunc func()

type task interface {
	cancel() error
	start() error
	refresh()
	watch()
}

const (
	TaskKeepLive TaskType = "KeepLive"
)

var taskList = make(map[string]map[TaskType]task)

func StopTask(deviceId string, taskType TaskType) error {
	task, err := getTask(deviceId, taskType)

	if err != nil {
		return err
	}

	delete(taskList[deviceId], taskType)

	return task.cancel()
}

func StartTask(deviceId string, taskType TaskType, duration time.Duration, runFunc runFunc) error {
	if taskList[deviceId] == nil {
		taskList[deviceId] = make(map[TaskType]task)
	}

	if taskList[deviceId][taskType] != nil {
		log.Errorf("task %+v already exists!", taskType)
		return errors.New(fmt.Sprintf("Start task error with task type: %+v, deviceId: %+v", taskType, deviceId))
	}

	task, err := createTask(deviceId, taskType, duration, runFunc)

	if err != nil {
		return err
	}

	return task.start()
}

func ResetTime(deviceId string, taskType TaskType) error {
	task, err := getTask(deviceId, taskType)

	if err != nil {
		return err
	}

	task.refresh()

	return nil
}

func getTask(deviceId string, taskType TaskType) (task, error) {
	if taskList[deviceId] == nil {
		log.Errorf("task %+v for deviceId: %+v not exists!", taskType, deviceId)
		return nil, errors.New(fmt.Sprintf("Stop task error with task type: %+v, deviceId: %+v", taskType, deviceId))
	}

	if taskList[deviceId][taskType] == nil {
		log.Errorf("task %+v for deviceId: %+v not exists!", taskType, deviceId)
		return nil, errors.New(fmt.Sprintf("Stop task error with task type: %+v, deviceId: %+v", taskType, deviceId))
	}

	task := taskList[deviceId][taskType]

	return task, nil
}

func createTask(deviceId string, taskType TaskType, duration time.Duration, runFunc runFunc) (task, error) {
	if taskList[deviceId][taskType] != nil {
		return nil, errors.New(fmt.Sprintf("error get task type object, task type: %+v", taskType))
	}
	var task task

	switch taskType {

	case TaskKeepLive:
		task = &keepLiveTask{
			deviceId: deviceId,
			duration: duration,
			runFunc:  runFunc,
		}
		taskList[deviceId][taskType] = task

	default:
		log.Errorf("Unsupported task type", taskType)
		return nil, errors.New(fmt.Sprintf("Unsupported task type with task type: %+v, deviceId: %+v",
			taskType, deviceId))
	}

	return task, nil
}
