package delay

import (
	"github.com/RussellLuo/timingwheel"
	"time"
	"xhyovo.cn/community/pkg/log"
)

var taskMap = make(map[int]*task)

type task struct {
	id    int // 任务标识
	timer *timingwheel.Timer
}

func newTask(id int, timer *timingwheel.Timer) *task {
	return &task{id: id, timer: timer}
}

var delay *DelayQueue

type DelayQueue struct {
	tw *timingwheel.TimingWheel
}

func init() {

	tw1 := timingwheel.NewTimingWheel(time.Second, 60)
	tw1.Start()
	delay = &DelayQueue{tw: tw1}
}

func GetInstant() *DelayQueue {
	return delay
}

func (delay *DelayQueue) Add(id int, expireTIme time.Time, taskFunction func()) {
	log.Infof("延迟队列加入任务：%v,触发时间为：%v", taskFunction, expireTIme.Format("2006-01-02 15:04:05"))
	// 过期时间 - 当前时间
	delay.tw.AfterFunc(expireTIme.Sub(time.Now()), taskFunction)
}

func (delay *DelayQueue) Remove(id int) {
	removeTask(id)
}

func addTask(task *task) {
	taskMap[task.id] = task
}

func removeTask(id int) {
	task, b := taskMap[id]
	if b {
		task.timer.Stop()
	}
}
