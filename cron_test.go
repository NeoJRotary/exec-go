package exec

import (
	"testing"
	"time"
)

func TestCron_StartStop(t *testing.T) {
	i := 0
	cron := NewCron()
	cron.DoJob(func() { i++ }).EverySecond()
	err := cron.Start()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(4 * time.Second)
	cron.Stop()
	nowI := i
	if i < 3 || i > 4 {
		t.Fatal("should run 3 ~ 4 times, but it did", i)
	}
	time.Sleep(2 * time.Second)
	if i-nowI > 1 {
		t.Fatal("should stop")
	}
}

func TestCron_MultiJobs(t *testing.T) {
	i := 0
	j := 0
	k := 0
	cron := NewCron()
	job := cron.DoJob(func() { i++ })
	job.EverySecond()
	cron.DoJob(func() { j++ }).Every(29)
	cron.DoJob(func() { k++ }).EveryMinute()

	cron.Stopped = false
	location := time.Now().Location()
	for s := 29; s <= 89; s++ {
		t := time.Date(2018, 1, 1, 1, s/60, s%60, 0, location)
		cron.doJobs(&t)
	}
	cron.Stopped = true
	time.Sleep(time.Second)
	if i != 61 {
		t.Fatal("i should run 61 times, but it did", i)
	}
	if j != 2 {
		t.Fatal("j should run 2 times, but it did", j)
	}
	if k != 1 {
		t.Fatal("k should run 1 time, but it did", k)
	}
}

func TestCron_EveryHour(t *testing.T) {
	i := 0
	cron := NewCron()
	cron.DoJob(func() { i++ }).EveryHour()

	cron.Stopped = false
	location := time.Now().Location()
	for n := 10; n < 180; n++ {
		t := time.Date(2018, 1, 1, n/60, n, 0, 0, location)
		cron.doJobs(&t)
	}
	cron.Stopped = true
	time.Sleep(time.Second)

	if i != 2 {
		t.Fatal("i should run 2 times, but it did", i)
	}
}

func TestCron_EveryDay(t *testing.T) {
	i := 0
	cron := NewCron()
	cron.DoJob(func() { i++ }).EveryDay()

	cron.Stopped = false
	location := time.Now().Location()
	for n := 2; n < 24*3; n++ {
		t := time.Date(2018, 1, n/24, n, 0, 0, 0, location)
		cron.doJobs(&t)
	}
	cron.Stopped = true
	time.Sleep(time.Second)

	if i != 2 {
		t.Fatal("i should run 2 times, but it did", i)
	}
}

func TestCron_EveryMonth(t *testing.T) {
	i := 0
	cron := NewCron()
	cron.DoJob(func() { i++ }).EveryMonth()

	cron.Stopped = false
	location := time.Now().Location()
	for m := 1; m < 3; m++ {
		for d := 0; d < 24; d++ {
			t := time.Date(2018, time.Month(m), d, 0, 0, 0, 0, location)
			cron.doJobs(&t)
		}
	}
	cron.Stopped = true
	time.Sleep(time.Second)

	if i != 2 {
		t.Fatal("i should run 2 times, but it did", i)
	}
}

func TestCron_EveryYear(t *testing.T) {
	i := 0
	cron := NewCron()
	cron.DoJob(func() { i++ }).EveryYear()

	cron.Stopped = false
	location := time.Now().Location()
	for y := 2015; y < 2018; y++ {
		for m := 1; m <= 12; m++ {
			t := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, location)
			cron.doJobs(&t)
		}
	}
	cron.Stopped = true
	time.Sleep(time.Second)

	if i != 3 {
		t.Fatal("i should run 3 times, but it did", i)
	}
}

func TestCron_EveryWeekDay(t *testing.T) {
	i := 0
	cron := NewCron()
	cron.DoJob(func() { i++ }).EveryWeekDay(1, 1, 2, 3)

	cron.Stopped = false
	location := time.Now().Location()
	for d := 1; d <= 15; d++ {
		t := time.Date(2018, 1, d, 3, 2, 1, 0, location)
		cron.doJobs(&t)
	}
	cron.Stopped = true
	time.Sleep(time.Second)

	if i != 3 {
		t.Fatal("i should run 3 times, but it did", i)
	}
}

func TestCron_EveryWeek(t *testing.T) {
	i := 0
	cron := NewCron()
	cron.DoJob(func() { i++ }).EveryWeek()

	cron.Stopped = false
	location := time.Now().Location()
	for d := 7; d <= 21; d++ {
		t := time.Date(2018, 1, d, 0, 0, 0, 0, location)
		cron.doJobs(&t)
	}
	cron.Stopped = true
	time.Sleep(time.Second)

	if i != 3 {
		t.Fatal("i should run 3 times, but it did", i)
	}
}
