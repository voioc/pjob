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

// func (s *ServerGroupService) List(page, pageSize int, filters ...interface{}) ([]*model.ServerGroup, int64, error) {
// 	offset := (page - 1) * pageSize
// 	data := make([]*model.ServerGroup, 0)

// 	// query := model.GetDB()
// 	// var count int
// 	condition := " 1 = 1 "
// 	if len(filters) > 0 {
// 		for k := 0; k < len(filters); k += 2 {
// 			condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
// 		}
// 	}

// 	total, err := model.GetDB().Where(condition).Count(&model.ServerGroup{})
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	if err := model.GetDB().Where(condition).Limit(pageSize, offset).Find(&data); err != nil {
// 		return nil, 0, err
// 	}

// 	// 	// query := orm.NewOrm().QueryTable(TableName("task"))
// 	// 	// if len(filters) > 0 {
// 	// 	// 	l := len(filters)
// 	// 	// 	for k := 0; k < l; k += 2 {
// 	// 	// 		query = query.Filter(filters[k].(string), filters[k+1])
// 	// 	// 	}
// 	// 	// }

// 	return data, total, nil
// }

func (s *ServerGroupService) ServerGroupLists(authStr string, adminId int) (sgl map[int]string) {
	Filters := make([]interface{}, 0)
	Filters = append(Filters, "status = ", 1)
	if authStr != "0" && adminId != 1 {
		serverGroupIdsArr := strings.Split(authStr, ",")
		serverGroupIds := make([]int, 0)
		for _, v := range serverGroupIdsArr {
			id, _ := strconv.Atoi(v)
			serverGroupIds = append(serverGroupIds, id)
		}
		Filters = append(Filters, "id", serverGroupIds)
	}

	// groupResult, n, _ := s.List(1, 1000, Filters...)
	group := make([]model.ServerGroup, 0)
	if err := model.List(&group, 1, 1000, Filters...); err != nil {
		fmt.Println(err.Error())
	}

	sgl = make(map[int]string, 0)
	for _, gv := range group {
		sgl[gv.ID] = gv.GroupName
	}

	//sgl[0] = "默认分组"
	return sgl
}

// 根据任务组id获取对应的名字
func (s *ServerGroupService) GroupIDName(ids string) (map[int]string, error) {
	ids = strings.Trim(strings.Trim(ids, ","), "")
	gid := strings.Split(ids, ",")
	fmt.Println(gid)

	group := make([]*model.ServerGroup, 0)
	// err := model.GetDB().Where("status = 1").In("id", gid).Find(&group)
	err := model.GetDB().Where("status = 1").Find(&group)
	if err != nil {
		return nil, err
	}

	data := map[int]string{}
	for _, gv := range group {
		data[gv.ID] = gv.GroupName
	}
	return data, nil
}

// func (s *ServerGroupService) Add(sg *model.ServerGroup) (int, error) {
// 	_, err := model.GetDB().Insert(sg)
// 	return sg.ID, err
// }

// func (s *ServerGroupService) Update(data *model.ServerGroup, args ...bool) error {
// 	if data.GroupName == "" {
// 		return fmt.Errorf("组名不能为空")
// 	}

// 	if len(args) > 0 && args[0] {
// 		if _, err := model.GetDB().ID(data.ID).Cols("status").Update(data); err != nil {
// 			return err
// 		}
// 	} else {
// 		if _, err := model.GetDB().Where("id = ?", data.ID).Update(data); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
