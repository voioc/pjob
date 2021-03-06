/************************************************************
** @Description: task
** @Author: george hao
** @Date:   2019-06-24 13:22
** @Last Modified by:  george hao
** @Last Modified time: 2019-06-24 13:22
*************************************************************/
package server

import (
	"github.com/astaxie/beego/logs"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/worker"
)

type RpcTask struct {
}

type RpcResult struct {
	Status  int
	Message string
}

//Execute once
func (r *RpcTask) RunTask(task *model.Task, Result *worker.JobResult) error {
	server_id := C.ServerId
	job, err := RestJobFromTask(task, server_id)
	if err != nil {
		return nil
	}
	*Result = *(Run(job))
	return nil
}

//Kill execution
func (r *RpcTask) KillCommand(task model.Task, reply *RpcResult) error {
	reply.Status = 200
	reply.Message = "Ok kill " + task.TaskName
	return nil
}

func (r *RpcTask) HeartBeat(ping string, reply *RpcResult) error {
	reply.Status = 200
	reply.Message = ping + " pong"
	logs.Info(ping)
	return nil
}
