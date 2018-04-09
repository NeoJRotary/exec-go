package exec

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Cron cron job scheduler
type Cron struct {
	Stopped bool
	jobI    uint32
	jobs    *sync.Map
}

// Job cron job instance
type Job struct {
	ID       uint32
	Stopped  bool
	cr       *Cron
	f        func()
	triggerI uint32
	triggers *sync.Map
}

// jobTrigger trigger of cron job
type jobTrigger struct {
	id         uint32
	lastTime   string
	weekday    string
	timePrefix string
}

// NewCron get new Cron obejct
func NewCron() *Cron {
	// time.Mo
	return &Cron{
		Stopped: true,
		jobs:    &sync.Map{},
	}
}

// DoJob set schedule by timeSyntax which is the same as crontab's
func (cr *Cron) DoJob(f func()) *Job {
	i := atomic.AddUint32(&cr.jobI, 1)
	j := &Job{
		ID:       i,
		cr:       cr,
		f:        f,
		triggers: &sync.Map{},
	}
	cr.jobs.Store(i, j)
	return j
}

// RemoveJob remove job by id
func (cr *Cron) RemoveJob(id uint32) {
	cr.jobs.Delete(id)
}

// get current time string
func getCurrent(now *time.Time) string {
	if now == nil {
		t := time.Now()
		now = &t
	}
	s := "00"
	s += fmt.Sprintf("%02d", now.Second())
	s += fmt.Sprintf("%02d", now.Minute())
	s += fmt.Sprintf("%02d", now.Hour())
	s += fmt.Sprintf("%02d", now.Day())
	s += fmt.Sprintf("%02d", now.Month())
	s += fmt.Sprintf("%02d", now.Weekday())
	return s
}

// run cron jobs
func (cr *Cron) doJobs(now *time.Time) {
	current := getCurrent(now)
	cr.jobs.Range(func(k, v interface{}) bool {
		j := v.(*Job)
		j.run(current)
		return !cr.Stopped
	})
}

// run cron jobs
func (cr *Cron) run() {
	for !cr.Stopped {
		time.Sleep(time.Second)
		if cr.Stopped {
			break
		}
		cr.doJobs(nil)
	}
}

// Start start cron jobs
func (cr *Cron) Start() error {
	if !cr.Stopped {
		return errors.New("Cron is running")
	}
	cr.Stopped = false
	go cr.run()
	return nil
}

// Stop stop cron jobs
func (cr *Cron) Stop() {
	cr.Stopped = true
}

// Stop stop this job
func (j *Job) Stop() {
	j.Stopped = true
}

// Restart restart this job
func (j *Job) Restart() {
	j.Stopped = false
}

func (j *Job) run(current string) {
	if j.Stopped || j.cr.Stopped {
		return
	}
	j.triggers.Range(func(k, v interface{}) bool {
		trg := v.(*jobTrigger)
		if !strings.HasPrefix(current, trg.timePrefix) {
			return !j.Stopped && !j.cr.Stopped
		} else if trg.weekday != "99" {
			if trg.weekday != current[12:14] {
				return !j.Stopped && !j.cr.Stopped
			}
		}
		go j.f()
		return false
	})
}

// addTrigger
func (j *Job) addTrigger(s string) uint32 {
	i := atomic.AddUint32(&j.triggerI, 1)
	t := &jobTrigger{
		id:         i,
		weekday:    "99",
		timePrefix: s,
	}
	j.triggers.Store(i, t)
	return i
}

// RemoveTrigger remove trigger by id
func (j *Job) RemoveTrigger(id uint32) {
	j.triggers.Delete(id)
}

// Every set specific time in [second(0~59), minute(0~59), hour(0~23), day(1~31), month(1~12)], unset part will be pretended as every (*). Invalid value will be set to first valid value. Return trigger id.
func (j *Job) Every(v ...int) uint32 {
	s := "00"
	for i, n := range v {
		switch i {
		case 0:
			fallthrough
		case 1:
			if n < 0 || n > 59 {
				n = 0
			}
		case 2:
			if n < 0 || n > 23 {
				n = 0
			}
		case 3:
			if n < 1 || n > 31 {
				n = 1
			}
		case 4:
			if n < 1 || n > 12 {
				n = 1
			}
		}
		s += fmt.Sprintf("%02d", n)
	}
	return j.addTrigger(s)
}

// EverySecond setup trigger which run job every second (**:**:**)
func (j *Job) EverySecond() uint32 {
	return j.Every()
}

// EveryMinute setup trigger which run job every minute (**:**:00)
func (j *Job) EveryMinute() uint32 {
	return j.Every(0)
}

// EveryHour setup trigger which run job every hour (**:00:00)
func (j *Job) EveryHour() uint32 {
	return j.Every(0, 0)
}

// EveryDay setup trigger which run job every day (00:00:00)
func (j *Job) EveryDay() uint32 {
	return j.Every(0, 0, 0)
}

// EveryMonth setup trigger which run job every month (1st 00:00:00)
func (j *Job) EveryMonth() uint32 {
	return j.Every(0, 0, 0, 1)
}

// EveryYear setup trigger which run job every year (Jan 1st 00:00:00)
func (j *Job) EveryYear() uint32 {
	return j.Every(0, 0, 0, 1, 1)
}

// EveryWeekDay setup trigger which run job every weekday at time. If time unset it will be 00:00:00.
func (j *Job) EveryWeekDay(weekday int, n ...int) uint32 {
	tid := j.Every(n...)
	v, _ := j.triggers.Load(tid)
	if weekday < 0 || weekday > 6 {
		weekday = 0
	}
	v.(*jobTrigger).weekday = fmt.Sprintf("%02d", weekday)
	return tid
}

// EveryWeek setup trigger which run job every week (Sunday 00:00:00)
func (j *Job) EveryWeek() uint32 {
	return j.EveryWeekDay(0, 0, 0, 0)
}
