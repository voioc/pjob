package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/notify"
	"github.com/voioc/cjob/worker"
)

type NotifyService struct {
	common.Base
}

// NotifyS instance
func NotifyS(c *gin.Context) *NotifyService {
	return &NotifyService{Base: common.Base{C: c}}
}

func (s *NotifyService) NotifyTypeList(ntype int) ([]*model.NotifyTpl, error) {
	data := make([]*model.NotifyTpl, 0)
	if err := model.GetDB().Where("tpl_type = ? and status = 1", ntype).Find(&data); err != nil {
		return nil, err

	}
	return data, nil
}

func (s *NotifyService) NotifyListIDs(ids []int) ([]*model.NotifyTpl, error) {
	data := make([]*model.NotifyTpl, 0)
	if err := model.GetDB().Where("status = 1").In("id", ids).Find(&data); err != nil {
		return nil, err

	}
	return data, nil
}

func (s *NotifyService) NotifyList(page, pageSize int, filters ...interface{}) ([]*model.NotifyTpl, int64, error) {
	offset := (page - 1) * pageSize
	data := make([]*model.NotifyTpl, 0)

	db := model.GetDB().Where("1=1")

	in := map[string]interface{}{}
	condition := " 1 = 1 "
	if len(filters) > 0 {
		for k := 0; k < len(filters); k += 2 {
			// 如果是数组则单独筛出来
			if _, flag := filters[k+1].([]int); flag {
				in[filters[k].(string)] = filters[k+1]
			} else {
				condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
			}
		}
	}

	if len(in) > 0 {
		for col, v := range in {
			if col != "" {
				regex := strings.Split(col, " ")
				if len(regex) == 2 && regex[1] == "not" {
					db = db.NotIn(col, v)
				} else {
					db = db.In(col, v)
				}
			}
		}
	}

	total, err := db.Where(condition).Count(&model.NotifyTpl{})
	if err != nil {
		return nil, 0, err
	}

	if err := db.Where(condition).Limit(pageSize, offset).Find(&data); err != nil {
		return nil, 0, err
	}

	// query := orm.NewOrm().QueryTable(TableName("task"))
	// if len(filters) > 0 {
	// 	l := len(filters)
	// 	for k := 0; k < l; k += 2 {
	// 		query = query.Filter(filters[k].(string), filters[k+1])
	// 	}
	// }

	return data, total, nil
}

// func (s *NotifyService) NotifyByID(id int) (*model.NotifyTpl, error) {
// 	data := &model.NotifyTpl{}

// 	if _, err := model.GetDB().Where("id = ?", id).Get(data); err != nil {
// 		return nil, err
// 	}

// 	if data.ID == 0 {
// 		return nil, fmt.Errorf("record not found")
// 	}

// 	return data, nil
// }

// func (s *NotifyService) Add(data *model.NotifyTpl) (int, error) {
// 	_, err := model.GetDB().Insert(data)
// 	return data.ID, err
// }

// func (s *NotifyService) Update(data *model.NotifyTpl, args ...bool) error {
// 	if len(args) > 0 && args[0] {
// 		if _, err := model.GetDB().Cols("status").Where("id = ?", data.ID).Update(data); err != nil {
// 			return err
// 		}
// 	} else {
// 		if _, err := model.GetDB().Where("id = ?", data.ID).Update(data); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *NotifyService) Del(ids interface{}) error {
// 	_, flag1 := ids.([]int)
// 	_, flag2 := ids.([]string)

// 	if flag1 || flag2 {
// 		_, err := model.GetDB().In("id", ids).Delete(&model.NotifyTpl{})
// 		return err
// 	}

// 	return nil
// }

func (s *NotifyService) NotifyFunc(job *worker.Job, result *worker.JobResult) int {
	if result.IsTimeout || !result.IsOk {
		task := model.Task{}
		if err := model.DataByID(&task, job.TaskID); err != nil {
			fmt.Println(err.Error())
			return 1
		}

		if task.IsNotify == 1 && task.NotifyUserIDs != "0" && task.NotifyUserIDs != "" {
			admin := make([]model.Admin, 0)
			// adminInfo := AllAdminInfo(j.Task.NotifyUserIds)
			if err := model.DataByIDs(&admin, strings.Split(task.NotifyUserIDs, ",")); err != nil {
				fmt.Println(err.Error())
			}

			if len(admin) == 0 {
				fmt.Println("no notify user")
			}

			phone := make(map[string]string, 0)
			dingtalk := make(map[string]string, 0)
			wechat := make(map[string]string, 0)
			toEmail := ""
			for _, v := range admin {
				if v.Phone != "0" && v.Phone != "" {
					phone[v.Phone] = v.Phone
				}

				if v.Email != "0" && v.Email != "" {
					toEmail += v.Email + ";"
				}

				if v.Dingtalk != "0" && v.Dingtalk != "" {
					dingtalk[v.Dingtalk] = v.Dingtalk
				}

				if v.Wechat != "0" && v.Wechat != "" {
					wechat[v.Wechat] = v.Wechat
				}
			}

			toEmail = strings.TrimRight(toEmail, ";")

			TextStatus := []string{
				"正常",
				"超时",
				"错误",
			}
			// status := log.Status + 2

			status := 0
			if result.IsTimeout {
				status = 1
			}

			if !result.IsOk {
				status = 2
			}

			title, content, taskOutput, errOutput := "", "", "", ""

			// notifyTpl, err := model.NotifyTplGetById(j.Task.NotifyTplID)
			notifyTpl := model.NotifyTpl{}
			if err := model.DataByID(&notifyTpl, task.NotifyTplID); err != nil {
				// notifyTpl, err := model.NotifyTplGetByTplType(task.NotifyType, model.NotifyTplTypeSystem)
				if flag, err := model.GetDB().Where("type = ? and tpl_type = ?", task.NotifyType, model.NotifyTplTypeSystem).Get(&notifyTpl); !flag || err != nil {
					msg := "record not found"
					if err != nil {
						msg = err.Error()
					}
					fmt.Println(msg)
					return 1
				}
			}

			title = notifyTpl.Title
			content = notifyTpl.Content

			taskOutput = strings.Replace(result.OutMsg, "\n", " ", -1)
			taskOutput = strings.Replace(taskOutput, "\"", "\\\"", -1)
			errOutput = strings.Replace(result.ErrMsg, "\n", " ", -1)
			errOutput = strings.Replace(errOutput, "\"", "\\\"", -1)

			if title != "" {
				title = strings.Replace(title, "{{TaskID}}", strconv.Itoa(job.TaskID), -1)
				title = strings.Replace(title, "{{ServerID}}", strconv.Itoa(job.ServerID), -1)
				title = strings.Replace(title, "{{TaskName}}", task.TaskName, -1)
				title = strings.Replace(title, "{{ExecuteCommand}}", task.Command, -1)
				title = strings.Replace(title, "{{ExecuteTime}}", job.StartAt.Format("2006-01-02 15:04:05"), -1)
				title = strings.Replace(title, "{{ProcessTime}}", strconv.FormatFloat(float64(int(job.ProcessTime))/1000, 'f', 6, 64), -1)
				title = strings.Replace(title, "{{ExecuteStatus}}", TextStatus[status], -1)
				title = strings.Replace(title, "{{TaskOutput}}", taskOutput, -1)
				title = strings.Replace(title, "{{ErrorOutput}}", errOutput, -1)
			}

			if content != "" {
				content = strings.Replace(content, "{{TaskID}}", strconv.Itoa(job.TaskID), -1)
				content = strings.Replace(content, "{{ServerID}}", strconv.Itoa(job.ServerID), -1)
				content = strings.Replace(content, "{{TaskName}}", task.TaskName, -1)
				content = strings.Replace(content, "{{ExecuteCommand}}", strings.Replace(task.Command, "\"", "\\\"", -1), -1)
				content = strings.Replace(content, "{{ExecuteTime}}", job.StartAt.Format("2006-01-02 15:04:05"), -1)
				content = strings.Replace(content, "{{ProcessTime}}", strconv.FormatFloat(float64(int(job.ProcessTime))/1000, 'f', 6, 64), -1)
				content = strings.Replace(content, "{{ExecuteStatus}}", TextStatus[status], -1)
				content = strings.Replace(content, "{{TaskOutput}}", taskOutput, -1)
				content = strings.Replace(content, "{{ErrorOutput}}", errOutput, -1)
			}

			if task.NotifyType == 1 && toEmail != "" {
				// 邮件
				mailtype := "html"
				ok := notify.SendToChan(toEmail, title, content, mailtype)
				if !ok {
					fmt.Println("发送邮件错误", toEmail)
				}
			} else if task.NotifyType == 2 && len(phone) > 0 {
				// 信息
				param := make(map[string]string)
				err := json.Unmarshal([]byte(content), &param)
				if err != nil {
					fmt.Println("发送信息错误", err)
					return 1
				}

				ok := notify.SendSmsToChan(phone, param)
				if !ok {
					fmt.Println("发送信息错误", phone)
				}
			} else if task.NotifyType == 3 && len(dingtalk) > 0 {
				//钉钉
				param := make(map[string]interface{})

				err := json.Unmarshal([]byte(content), &param)
				if err != nil {
					fmt.Println("发送钉钉错误", err)
					return 1
				}

				ok := notify.SendDingtalkToChan(dingtalk, param)
				if !ok {
					fmt.Println("发送钉钉错误", dingtalk)
				}
			} else if task.NotifyType == 4 && len(wechat) > 0 {
				//微信
				param := make(map[string]string)
				err := json.Unmarshal([]byte(content), &param)
				if err != nil {
					fmt.Println("发送微信错误", err)
					return 1
				}

				ok := notify.SendWechatToChan(phone, param)
				if !ok {
					fmt.Println("发送微信错误", phone)
				}
			}
		}
	}

	return 1
}
