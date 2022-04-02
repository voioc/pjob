package service

import (
	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/define"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
)

type ServerService struct {
	common.Base
}

// ServerS instance
func ServerS(c *gin.Context) *ServerService {
	return &ServerService{Base: common.Base{C: c}}
}

func (s *ServerService) ServersListID(ids interface{}) ([]*model.TaskServer, error) {
	_, flag1 := ids.([]int)
	_, flag2 := ids.([]string)

	db := model.GetDB().Where("status = 1")
	if flag1 || flag2 {
		db = db.In("id", ids)
	}

	data := make([]*model.TaskServer, 0)
	if err := db.Find(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// // 根据任务组id获取对应的名字
// func (s *ServerService) GroupIDName(ids interface{}) (map[int]string, error) {
// 	_, flag1 := ids.([]int)
// 	_, flag2 := ids.([]string)

// 	group := make([]*model.TaskServer, 0)
// 	db := model.GetDB().Where("status = 1")
// 	if flag1 || flag2 {
// 		db = db.In("id", ids)
// 	}
// 	// err := model.GetDB().Where("status = 1").In("id", gid).Find(&group)
// 	err := model.GetDB().Find(&group)
// 	if err != nil {
// 		return nil, err
// 	}

// 	data := map[int]string{}
// 	for _, gv := range group {
// 		data[gv.ID] = gv.ServerName
// 	}
// 	return data, nil
// }

func (s *ServerService) ServerLists(ServerGroupIDS string) ([]define.ServerList, error) {
	// 获取有权限的用户组
	serverGroup, err := ServerGroupS(s.C).GroupIDName(ServerGroupIDS)
	if err != nil {
		return nil, err
	}

	// 获取所有服务器
	serverList := make([]*model.TaskServer, 0)
	if err := model.GetDB().In("status", []int{0, 1}).Find(&serverList); err != nil {
		return nil, err
	}

	data := make([]define.ServerList, 0)
	for k, v := range serverGroup {
		sl := define.ServerList{}
		sl.GroupID = k
		sl.GroupName = v
		servers := make(map[int]string)

		for _, sv := range serverList {
			if sv.GroupID == k {
				servers[sv.ID] = sv.ServerName
			}
		}

		sl.Servers = servers
		data = append(data, sl)
	}
	return data, nil
}
