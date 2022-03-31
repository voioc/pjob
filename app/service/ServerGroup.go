package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
)

type ServerGroupService struct {
	common.Base
}

// ServerGroupS instance
func ServerGroupS(c *gin.Context) *ServerGroupService {
	return &ServerGroupService{Base: common.Base{C: c}}
}

func (s *ServerGroupService) List(page, pageSize int, filters ...interface{}) ([]*model.ServerGroup, int64, error) {
	offset := (page - 1) * pageSize
	data := make([]*model.ServerGroup, 0)

	// query := model.GetDB()
	// var count int
	condition := " 1 = 1 "
	if len(filters) > 0 {
		for k := 0; k < len(filters); k += 2 {
			condition = fmt.Sprintf("%s and %s %s", condition, filters[k].(string), filters[k+1])
		}
	}

	total, err := model.GetDB().Where(condition).Count(&model.ServerGroup{})
	if err != nil {
		return nil, 0, err
	}

	if err := model.GetDB().Where(condition).Limit(pageSize, offset).Find(&data); err != nil {
		return nil, 0, err
	}

	// 	// query := orm.NewOrm().QueryTable(TableName("task"))
	// 	// if len(filters) > 0 {
	// 	// 	l := len(filters)
	// 	// 	for k := 0; k < l; k += 2 {
	// 	// 		query = query.Filter(filters[k].(string), filters[k+1])
	// 	// 	}
	// 	// }

	return data, total, nil
}

func (s *ServerGroupService) ServerGroupLists(authStr string, adminId int) (sgl map[int]string) {
	Filters := make([]interface{}, 0)
	Filters = append(Filters, "status", 1)
	if authStr != "0" && adminId != 1 {
		serverGroupIdsArr := strings.Split(authStr, ",")
		serverGroupIds := make([]int, 0)
		for _, v := range serverGroupIdsArr {
			id, _ := strconv.Atoi(v)
			serverGroupIds = append(serverGroupIds, id)
		}
		Filters = append(Filters, "id__in", serverGroupIds)
	}

	groupResult, n, _ := s.List(1, 1000, Filters...)

	sgl = make(map[int]string, n)
	for _, gv := range groupResult {
		sgl[gv.ID] = gv.GroupName
	}

	//sgl[0] = "默认分组"
	return sgl
}

// func (s *TaskLogService) GetLogNum(status int) (int64, error) {
// 	// return orm.NewOrm().QueryTable(TableName("task_log")).Filter("status", status).Count()

// 	return model.GetDB().Where("status = ?", status).Count(&model.TaskLog{})
// }
