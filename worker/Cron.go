package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/voioc/coco/logzap"
)

var (
	mainCron *cron.Cron
	workPool chan bool
	lock     sync.Mutex
	pool     map[int]cron.EntryID
)

func init() {
	if size := viper.GetInt("job.pool"); size > 0 {
		workPool = make(chan bool, size)
	}

	pool = map[int]cron.EntryID{}
	mainCron = cron.New(cron.WithSeconds())
	mainCron.Start()
}

func AddJob(spec string, job *Job) bool {
	lock.Lock()
	defer lock.Unlock()

	if entry, _ := GetEntryByID(job.JobKey); entry != nil {
		return false
	}

	id, err := mainCron.AddJob(spec, job)
	if err != nil {
		logzap.Ex(context.Background(), "Cron ", "AddJob error: %s", err.Error())
		return false
	}

	pool[job.JobKey] = id
	// fmt.Println("id: ", id)
	return true
}

func RemoveJob(jobKey int) error {
	entryID, flag := pool[jobKey]
	if !flag {
		logzap.Ex(context.Background(), "Cron", "DeleteJob not found: %d", jobKey)
		return fmt.Errorf("record not found: %d", jobKey)
	}

	mainCron.Remove(entryID)
	delete(pool, jobKey)
	return nil
}

func GetEntryByID(jobKey int) (*cron.Entry, error) {
	entryID, flag := pool[jobKey]
	if !flag {
		logzap.Ex(context.Background(), "Cron", "DeleteJob not found: %d", jobKey)
		return nil, fmt.Errorf("record not found: %d", jobKey)
	}

	e := mainCron.Entry(entryID)
	return &e, nil
}

func GetEntries(size int) []cron.Entry {
	ret := mainCron.Entries()
	if len(ret) > size {
		return ret[:size]
	}

	return ret
}
