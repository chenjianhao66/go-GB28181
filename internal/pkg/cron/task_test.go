package cron

import (
	"github.com/agiledragon/gomonkey"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

const fakeDeviceId = "400001000"
const fakeDuration = 3 * time.Second

var fakeTaskList taskSchedule = make(map[string]map[TaskType]task)

func replaceAndGenerateGlobalVar() task {
	fakeTaskList[fakeDeviceId] = make(map[TaskType]task)
	fakeKeepLiveTask := &keepLiveTask{
		deviceId: fakeDeviceId,
		duration: fakeDuration,
		runFunc: func() {
			return
		},
	}
	fakeTaskList[fakeDeviceId][TaskKeepLive] = fakeKeepLiveTask
	gomonkey.ApplyGlobalVar(&taskList, fakeTaskList)
	return fakeKeepLiveTask
}

func TestResetTime(t *testing.T) {
	convey.Convey("TestResetTime", t, func() {
		convey.Convey("for success", func() {
			fakeKeepLiveTask := replaceAndGenerateGlobalVar()
			patches := gomonkey.ApplyFunc(getTask, func(deviceId string, taskType TaskType) (task, error) {
				return fakeKeepLiveTask, nil
			})

			err := fakeKeepLiveTask.start()
			convey.So(err, convey.ShouldEqual, nil)
			convey.So(fakeTaskList.getOneTask(fakeDeviceId, TaskKeepLive), convey.ShouldNotEqual, nil)

			defer patches.Reset()

			err = ResetTime(fakeDeviceId, TaskKeepLive)
			convey.So(err, convey.ShouldEqual, nil)

			time.Sleep(5 * time.Second)

			convey.So(fakeTaskList.getOneTask(fakeDeviceId, TaskKeepLive), convey.ShouldEqual, nil)
		})
	})
}

func TestStartTask(t *testing.T) {
	convey.Convey("TestStartTask", t, func() {
		convey.Convey("for success", func() {
			gomonkey.ApplyGlobalVar(&taskList, fakeTaskList)
			err := StartTask(fakeDeviceId, TaskKeepLive, fakeDuration, func() {})
			convey.So(err, convey.ShouldEqual, nil)
		})
	})
}

func TestStopTask(t *testing.T) {
	convey.Convey("TestStopTask", t, func() {
		convey.Convey("for success", func() {
			tt := replaceAndGenerateGlobalVar()

			fakeTask := tt.(*keepLiveTask)
			err := fakeTask.start()
			convey.So(err, convey.ShouldEqual, nil)

			err = StopTask(fakeDeviceId, TaskKeepLive)
			convey.So(err, convey.ShouldEqual, nil)
		})
	})
}

func Test_createTask(t *testing.T) {
	convey.Convey("Test_createTask", t, func() {
		convey.Convey("for success", func() {
			gomonkey.ApplyGlobalVar(&taskList, fakeTaskList)
			t2, err := createTask(fakeDeviceId, TaskKeepLive, fakeDuration, func() {})
			convey.So(t2, convey.ShouldNotEqual, nil)

			liveTask := t2.(*keepLiveTask)
			convey.So(liveTask.deviceId, convey.ShouldEqual, fakeDeviceId)
			convey.So(liveTask.timer, convey.ShouldEqual, nil)
			convey.So(liveTask.duration, convey.ShouldEqual, fakeDuration)

			convey.So(err, convey.ShouldEqual, nil)
		})
	})
}

func Test_getTask(t *testing.T) {
	convey.Convey("Test_getTask", t, func() {
		convey.Convey("for success", func() {
			replaceAndGenerateGlobalVar()
			t2, err := getTask(fakeDeviceId, TaskKeepLive)
			liveTask := t2.(*keepLiveTask)
			convey.So(liveTask.deviceId, convey.ShouldEqual, fakeDeviceId)
			convey.So(liveTask.timer, convey.ShouldEqual, nil)
			convey.So(liveTask.duration, convey.ShouldEqual, fakeDuration)

			convey.So(err, convey.ShouldEqual, nil)
		})
	})
}

func Test_taskSchedule_clearAllTasksForOneDevice(t *testing.T) {
	convey.Convey("Test_taskSchedule_clearAllTasksForOneDevice", t, func() {
		convey.Convey("for success", func() {
			replaceAndGenerateGlobalVar()
			_, err := getTask(fakeDeviceId, TaskKeepLive)

			convey.So(err, convey.ShouldEqual, nil)

			fakeTaskList.clearAllTasksForOneDevice(fakeDeviceId)

			convey.So(fakeTaskList[fakeDeviceId], convey.ShouldEqual, nil)
		})
	})
}

func Test_taskSchedule_deleteOneDeviceRecord(t *testing.T) {
	convey.Convey("Test_taskSchedule_deleteOneDeviceRecord", t, func() {
		convey.Convey("for success", func() {
			replaceAndGenerateGlobalVar()
			_, err := getTask(fakeDeviceId, TaskKeepLive)

			convey.So(err, convey.ShouldEqual, nil)

			fakeTaskList.deleteOneDeviceRecord(fakeDeviceId)

			convey.So(fakeTaskList[fakeDeviceId], convey.ShouldEqual, nil)
		})
	})
}

func Test_taskSchedule_deleteOneTask(t *testing.T) {
	convey.Convey("Test_taskSchedule_deleteOneTask", t, func() {
		convey.Convey("for success", func() {
			replaceAndGenerateGlobalVar()
			_, err := getTask(fakeDeviceId, TaskKeepLive)

			convey.So(err, convey.ShouldEqual, nil)

			fakeTaskList.deleteOneTask(fakeDeviceId, TaskKeepLive)

			convey.So(fakeTaskList[fakeDeviceId], convey.ShouldNotEqual, nil)
			convey.So(fakeTaskList[fakeDeviceId][TaskKeepLive], convey.ShouldEqual, nil)
		})
	})
}

func Test_taskSchedule_getAllTasksForOneDevice(t *testing.T) {
	convey.Convey("Test_taskSchedule_getAllTasksForOneDevice", t, func() {
		convey.Convey("for success", func() {
			replaceAndGenerateGlobalVar()

			allTasks := fakeTaskList.getAllTasksForOneDevice(fakeDeviceId)

			convey.So(allTasks[TaskKeepLive], convey.ShouldNotEqual, nil)

			liveTask := allTasks[TaskKeepLive].(*keepLiveTask)
			convey.So(liveTask.deviceId, convey.ShouldEqual, fakeDeviceId)
			convey.So(liveTask.timer, convey.ShouldEqual, nil)
			convey.So(liveTask.duration, convey.ShouldEqual, fakeDuration)

			// 不存在的任务类型，空
			convey.So(allTasks["fake"], convey.ShouldEqual, nil)
		})
	})
}

func Test_taskSchedule_getOneTask(t *testing.T) {
	convey.Convey("Test_taskSchedule_getOneTask", t, func() {
		convey.Convey("for success", func() {
			replaceAndGenerateGlobalVar()

			oneTask := fakeTaskList.getOneTask(fakeDeviceId, TaskKeepLive)

			convey.So(oneTask, convey.ShouldNotEqual, nil)

			liveTask := oneTask.(*keepLiveTask)
			convey.So(liveTask.deviceId, convey.ShouldEqual, fakeDeviceId)
			convey.So(liveTask.timer, convey.ShouldEqual, nil)
			convey.So(liveTask.duration, convey.ShouldEqual, fakeDuration)
		})
	})
}

func Test_taskSchedule_initTaskList(t *testing.T) {
	convey.Convey("Test_taskSchedule_initTaskList", t, func() {
		convey.Convey("for success", func() {
			gomonkey.ApplyGlobalVar(&taskList, fakeTaskList)
			fakeTaskList.initTaskList(fakeDeviceId)

			convey.So(fakeTaskList[fakeDeviceId], convey.ShouldNotEqual, nil)
		})
	})
}

func Test_taskSchedule_putOneTask(t *testing.T) {
	convey.Convey("Test_taskSchedule_putOneTask", t, func() {
		convey.Convey("for success", func() {
			gomonkey.ApplyGlobalVar(&taskList, fakeTaskList)
			fakeTaskList.initTaskList(fakeDeviceId)

			fakeKeepLiveTask := &keepLiveTask{
				deviceId: fakeDeviceId,
				duration: fakeDuration,
				runFunc: func() {
					return
				},
			}
			fakeTaskList.putOneTask(fakeDeviceId, TaskKeepLive, fakeKeepLiveTask)

			convey.So(fakeTaskList[fakeDeviceId], convey.ShouldNotEqual, nil)
			convey.So(fakeTaskList[fakeDeviceId][TaskKeepLive], convey.ShouldNotEqual, nil)

			fakeTask := fakeTaskList[fakeDeviceId][TaskKeepLive]
			t2 := fakeTask.(*keepLiveTask)

			convey.So(t2.deviceId, convey.ShouldEqual, fakeDeviceId)
			convey.So(t2.timer, convey.ShouldEqual, nil)
			convey.So(t2.duration, convey.ShouldEqual, fakeDuration)
		})
	})
}
