package cron

import (
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/pkg/errors"
	"sync"
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

type taskSchedule map[string]map[TaskType]task

// key: deviceId, value: {key: taskType, value: task Object}
var taskList taskSchedule = make(map[string]map[TaskType]task)

var once sync.Once

func (t taskSchedule) deleteOneTask(deviceId string, taskType TaskType) {
	delete(t[deviceId], taskType)
}

func (t taskSchedule) clearAllTasksForOneDevice(deviceId string) {
	t[deviceId] = nil
}

func (t taskSchedule) deleteOneDeviceRecord(deviceId string) {
	delete(t, deviceId)
}

func (t taskSchedule) getOneTask(deviceId string, taskType TaskType) task {
	if t[deviceId] != nil {
		return t[deviceId][taskType]
	}
	return nil
}

func (t taskSchedule) getAllTasksForOneDevice(deviceId string) map[TaskType]task {
	return t[deviceId]
}

func (t taskSchedule) putOneTask(deviceId string, taskType TaskType, oneTask task) {
	if t[deviceId] != nil && t[deviceId][taskType] == nil {
		t[deviceId][taskType] = oneTask
	}
}

func (t taskSchedule) initTaskList(deviceId string) {
	if taskList.getAllTasksForOneDevice(deviceId) == nil {
		taskList[deviceId] = make(map[TaskType]task)
	}
}

func StopTask(deviceId string, taskType TaskType) error {
	t, err := getTask(deviceId, taskType)

	if err != nil {
		return err
	}

	taskList.deleteOneTask(deviceId, taskType)

	return t.cancel()
}

func StartTask(deviceId string, taskType TaskType, duration time.Duration, runFunc runFunc) error {
	if taskList.getOneTask(deviceId, taskType) != nil {
		log.Errorf("task %+v already exists!", taskType)
		return errors.New(fmt.Sprintf("Start task error with task type: %+v, deviceId: %+v", taskType, deviceId))
	}

	t, err := createTask(deviceId, taskType, duration, runFunc)

	if err != nil {
		return err
	}

	return t.start()
}

func ResetTime(deviceId string, taskType TaskType) error {
	t, err := getTask(deviceId, taskType)

	if err != nil {
		return err
	}

	t.refresh()

	return nil
}

func getTask(deviceId string, taskType TaskType) (task, error) {
	if taskList.getAllTasksForOneDevice(deviceId) == nil {
		log.Errorf("task %+v for deviceId: %+v not exists!", taskType, deviceId)
		return nil, errors.New(fmt.Sprintf("Stop task error with task type: %+v, deviceId: %+v", taskType, deviceId))
	}

	if taskList.getOneTask(deviceId, taskType) == nil {
		log.Errorf("task %+v for deviceId: %+v not exists!", taskType, deviceId)
		return nil, errors.New(fmt.Sprintf("Stop task error with task type: %+v, deviceId: %+v", taskType, deviceId))
	}

	t := taskList.getOneTask(deviceId, taskType)

	return t, nil
}

func createTask(deviceId string, taskType TaskType, duration time.Duration, runFunc runFunc) (task, error) {
	if taskList.getOneTask(deviceId, taskType) != nil {
		return nil, errors.New(fmt.Sprintf("error get task type object, task type: %+v", taskType))
	}

	taskList.initTaskList(deviceId)

	var t task

	switch taskType {

	case TaskKeepLive:
		t = &keepLiveTask{
			deviceId: deviceId,
			duration: duration,
			runFunc:  runFunc,
		}
		taskList.putOneTask(deviceId, taskType, t)

	default:
		log.Errorf("Unsupported task type", taskType)
		return nil, errors.New(fmt.Sprintf("Unsupported task type with task type: %+v, deviceId: %+v",
			taskType, deviceId))
	}

	return t, nil
}