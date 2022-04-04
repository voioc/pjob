/*
* @Author: haodaquan
* @Date:   2017-06-21 12:55:19
* @Last Modified by:   haodaquan
* @Last Modified time: 2017-06-21 13:03:06
 */

package jobs

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/astaxie/beego"
	"github.com/voioc/cjob/app/model"
)

func InitJobs() {
	// list, _ := model.TaskGetList(1, 1000000, "status", 1)
	list := make([]*model.Task, 0)
	if err := model.List(&list, 1, 1000000, "status =", 1); err != nil {
		fmt.Println(err.Error())
	}

	for _, task := range list {
		jobs, err := NewJobFromTask(task)
		if err != nil {
			beego.Error("InitJobs:", err.Error())
			continue
		}

		for _, job := range jobs {
			AddJob(task.CronSpec, job)
		}

	}
}

func runCmdWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):
		beego.Warn(fmt.Sprintf("任务执行时间超过%d秒，进程将被强制杀掉: %d", int(timeout/time.Second), cmd.Process.Pid))
		go func() {
			<-done // 读出上面的goroutine数据，避免阻塞导致无法退出
		}()
		if err = cmd.Process.Kill(); err != nil {
			beego.Error(fmt.Sprintf("进程无法杀掉: %d, 错误信息: %s", cmd.Process.Pid, err))
		}
		return err, true
	case err = <-done:
		return err, false
	}
}