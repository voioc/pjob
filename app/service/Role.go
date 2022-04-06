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

// func (s *RoleService) RoleList(page, pageSize int, filters ...interface{}) ([]*model.Role, error) {
// 	offset := (page - 1) * pageSize
// 	data := make([]*model.Role, 0)

// 	in := map[string]interface{}{}
// 	condition := " 1 = 1 "
// 	if len(filters) > 0 {
// 		for k := 0; k < len(filters); k += 2 {
// 			// 如果是数组则单独筛出来
// 			if _, flag := filters[k+1].([]int); flag {
// 				in[filters[k].(string)] = filters[k+1]
// 			} else {
// 				condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
// 			}
// 		}
// 	}

// 	db := model.GetDB().Where("1=1")
// 	if len(in) > 0 {
// 		for col, v := range in {
// 			if col != "" {
// 				regex := strings.Split(col, " ")
// 				if len(regex) == 2 && regex[1] == "not" {
// 					db = db.NotIn(col, v)
// 				} else {
// 					db = db.In(col, v)
// 				}
// 			}
// 		}
// 	}

// 	// total, err := db.Where(condition).Count(&model.Role{})
// 	// if err != nil {
// 	// 	return nil, 0, err
// 	// }

// 	if err := db.Where(condition).Limit(pageSize, offset).Find(&data); err != nil {
// 		return nil, err
// 	}

// 	return data, nil
// }

// func (s *RoleService) RoleListCount(filters ...interface{}) (int64, error) {
// 	// offset := (page - 1) * pageSize
// 	// data := make([]*model.Role, 0)

// 	in := map[string]interface{}{}
// 	condition := " 1 = 1 "
// 	if len(filters) > 0 {
// 		for k := 0; k < len(filters); k += 2 {
// 			// 如果是数组则单独筛出来
// 			if _, flag := filters[k+1].([]int); flag {
// 				in[filters[k].(string)] = filters[k+1]
// 			} else {
// 				condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
// 			}
// 		}
// 	}

// 	db := model.GetDB().Where("1=1")
// 	if len(in) > 0 {
// 		for col, v := range in {
// 			if col != "" {
// 				regex := strings.Split(col, " ")
// 				if len(regex) == 2 && regex[1] == "not" {
// 					db = db.NotIn(col, v)
// 				} else {
// 					db = db.In(col, v)
// 				}
// 			}
// 		}
// 	}

// 	return db.Where(condition).Count(&model.Role{})
// }

func (s *RoleService) TaskGroups(uid int, roleIDs string) (string, string) {
	if uid == 1 || roleIDs == "0" {
		return "", ""
	}

	filters := make([]interface{}, 0)
	filters = append(filters, "status = ", 1)

	RoleIdsArr := strings.Split(roleIDs, ",")

	RoleIds := make([]int, 0)
	for _, v := range RoleIdsArr {
		id, _ := strconv.Atoi(v)
		RoleIds = append(RoleIds, id)
	}

	filters = append(filters, "id", RoleIds)

	// result, _ := s.RoleList(1, 1000, filters...)
	result := make([]model.Role, 0)
	if err := model.List(&result, 1, 1000, filters...); err != nil {
		fmt.Println(err.Error())
	}

	serverGroups := ""
	taskGroups := ""
	for _, v := range result {
		serverGroups += v.ServerGroupIDs + ","
		taskGroups += v.TaskGroupIDs + ","
	}

	return strings.Trim(serverGroups, ","), strings.Trim(taskGroups, ",")
}

func (s *RoleService) Resources(uid int, roleIDs string) ([]int, []int) {
	if uid == 1 || roleIDs == "0" {
		return []int{}, []int{}
	}

	filters := make([]interface{}, 0)
	filters = append(filters, "status = ", 1)

	RoleIdsArr := strings.Split(roleIDs, ",")

	RoleIds := make([]int, 0)
	for _, v := range RoleIdsArr {
		id, _ := strconv.Atoi(v)
		RoleIds = append(RoleIds, id)
	}

	filters = append(filters, "id", RoleIds)

	// result, _ := s.RoleList(1, 1000, filters...)
	result := make([]model.Role, 0)
	if err := model.List(&result, 1, 1000, filters...); err != nil {
		fmt.Println(err.Error())
	}

	serverGroups := []int{}
	taskGroups := []int{}
	for _, v := range result {
		for _, tid := range strings.Split(v.TaskGroupIDs, ",") {
			tidInt, _ := strconv.Atoi(tid)
			taskGroups = append(taskGroups, tidInt)
		}
		for _, sid := range strings.Split(v.ServerGroupIDs, ",") {
			sidInt, _ := strconv.Atoi(sid)
			serverGroups = append(serverGroups, sidInt)
		}

		// serverGroups += v.ServerGroupIDs + ","
		// taskGroups += v.TaskGroupIDs + ","
	}

	return taskGroups, serverGroups
}

// func (s *RoleService) RoleByID(id int) (*model.Role, error) {
// 	data := &model.Role{}

// 	if _, err := model.GetDB().Where("id = ?", id).Get(data); err != nil {
// 		return nil, err
// 	}

// 	if data.ID == 0 {
// 		return nil, fmt.Errorf("server not found")
// 	}

// 	return data, nil
// }

// func (s *RoleService) Add(data *model.Role) (int, error) {
// 	_, err := model.GetDB().Insert(data)
// 	return data.ID, err
// }

// func (s *RoleService) Update(data *model.Role, args ...bool) error {
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

// func (s *RoleService) Del(ids interface{}) error {
// 	_, flag1 := ids.([]int)
// 	_, flag2 := ids.([]string)

// 	if flag1 || flag2 {
// 		_, err := model.GetDB().In("id", ids).Delete(&model.Role{})
// 		return err
// 	}

// 	return nil
// }
