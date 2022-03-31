package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
)

type RoleService struct {
	common.Base
}

// RoleS instance
func RoleS(c *gin.Context) *RoleService {
	return &RoleService{Base: common.Base{C: c}}
}

func (s *RoleService) RoleList(page, pageSize int, filters ...interface{}) ([]*model.Role, int64, error) {
	offset := (page - 1) * pageSize
	data := make([]*model.Role, 0)

	db := model.GetDB().Where("1=1")

	in := map[string]interface{}{}
	condition := " 1 = 1 "
	if len(filters) > 0 {
		for k := 0; k < len(filters); k += 2 {
			// 如果是数组则单独筛出来
			if _, flag := filters[k+1].([]int); !flag {
				in[filters[k].(string)] = filters[k+1]
			}

			condition = fmt.Sprintf("%s and %s %s", condition, filters[k].(string), filters[k+1])
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

	total, err := db.Where(condition).Count(&model.Role{})
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

func (s *RoleService) TaskGroups(uid int, roleIDs string) (string, string) {
	if uid == 1 || roleIDs == "0" {
		return "", ""
	}

	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)

	RoleIdsArr := strings.Split(roleIDs, ",")

	RoleIds := make([]int, 0)
	for _, v := range RoleIdsArr {
		id, _ := strconv.Atoi(v)
		RoleIds = append(RoleIds, id)
	}

	filters = append(filters, "id__in", RoleIds)

	result, _, _ := s.RoleList(1, 1000, filters...)
	serverGroups := ""
	taskGroups := ""
	for _, v := range result {
		serverGroups += v.ServerGroupIDs + ","
		taskGroups += v.TaskGroupIDs + ","
	}

	return strings.Trim(serverGroups, ","), strings.Trim(taskGroups, ",")
}
