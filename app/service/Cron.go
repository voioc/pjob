package service

import (
	"sync"

	"github.com/astaxie/beego"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/voioc/cjob/app/jobs"
	"github.com/voioc/cjob/common"
	cron "github.com/voioc/cjob/crons"
)

var (
	mainCron *cron.Cron
	workPool chan bool
	lock     sync.Mutex
)

func init() {
	if size := viper.GetInt("job.size"); size > 0 {
		workPool = make(chan bool, size)
	}

	mainCron = cron.New()
	mainCron.Start()
}

type CronService struct {
	common.Base
}

// RoleS instance
func CronS(c *gin.Context) *CronService {
	return &CronService{Base: common.Base{C: c}}
}

func (s *CronService) AddJob(spec string, job *jobs.Job) bool {
	lock.Lock()
	defer lock.Unlock()

	if s.GetEntryByID(job.JobKey) != nil {
		return false
	}
	err := mainCron.AddJob(spec, job)
	if err != nil {
		beego.Error("AddJob: ", err.Error())
		return false
	}
	//fmt.Println(job)
	return true
}

func (s *CronService) RemoveJob(jobKey int) {
	mainCron.RemoveJob(func(e *cron.Entry) bool {
		if v, ok := e.Job.(*jobs.Job); ok {
			if v.JobKey == jobKey {
				return true
			}
		}
		return false
	})
}

func (s *CronService) GetEntryByID(jobKey int) *cron.Entry {
	entries := mainCron.Entries()
	for _, e := range entries {
		if v, ok := e.Job.(*jobs.Job); ok {
			if v.JobKey == jobKey {
				return e
			}
		}
	}
	return nil
}

func (s *CronService) GetEntries(size int) []*cron.Entry {
	ret := mainCron.Entries()
	if len(ret) > size {
		return ret[:size]
	}
	return ret
}
