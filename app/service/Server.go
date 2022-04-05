package service

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/define"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/libs"
)

type ServerService struct {
	common.Base
}

// ServerS instance
func ServerS(c *gin.Context) *ServerService {
	return &ServerService{Base: common.Base{C: c}}
}

// func (s *ServerService) ServerList(page, pageSize int, filters ...interface{}) ([]*model.TaskServer, int64, error) {
// 	offset := (page - 1) * pageSize
// 	data := make([]*model.TaskServer, 0)

// 	// db := model.GetDB().Where("1=1")
// 	// clone := model.GetDB().Where("1=1")

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

// 	db := model.GetDB().Where("1=1").Where(condition)
// 	clone := model.GetDB().Where("1=1").Where(condition)
// 	if len(in) > 0 {
// 		for col, v := range in {
// 			if col != "" {
// 				regex := strings.Split(col, " ")
// 				if len(regex) == 2 && regex[1] == "not" {
// 					db = db.NotIn(col, v)
// 					clone = db.NotIn(col, v)
// 				} else {
// 					db = db.In(col, v)
// 					clone = db.In(col, v)
// 				}
// 			}
// 		}
// 	}

// 	x := *db
// 	clone = &x
// 	fmt.Printf("%p \n", db)
// 	fmt.Printf("%p \n", clone)

// 	total, err := clone.Count(&model.TaskServer{})
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	if err := db.OrderBy("field(status, 1, 2, 3, 0), id desc ").Limit(pageSize, offset).Find(&data); err != nil {
// 		return nil, 0, err
// 	}

// 	// query := orm.NewOrm().QueryTable(TableName("task"))
// 	// if len(filters) > 0 {
// 	// 	l := len(filters)
// 	// 	for k := 0; k < l; k += 2 {
// 	// 		query = query.Filter(filters[k].(string), filters[k+1])
// 	// 	}
// 	// }

// 	return data, total, nil
// }

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

// func (s *ServerService) ServerByID(id int) (*model.TaskServer, error) {
// 	data := &model.TaskServer{}

// 	if _, err := model.GetDB().Where("id = ?", id).Get(data); err != nil {
// 		return nil, err
// 	}

// 	if data.ID == 0 {
// 		return nil, fmt.Errorf("server not found")
// 	}

// 	return data, nil
// }

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

// func (s *ServerService) Add(server *model.TaskServer) (int, error) {
// 	_, err := model.GetDB().Insert(server)
// 	return server.ID, err
// }

// func (s *ServerService) Update(data *model.TaskServer, args ...bool) error {
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

// 服务器探活服务
func (s *ServerService) Probe(sid int) bool {

	//判断执行器或者服务器是否存活
	// server, _ := model.TaskServerGetById(pollServerId)
	server := model.TaskServer{}
	if err := model.DataByID(&server, sid); err != nil {
		fmt.Println(err.Error())
		return false
	}

	if server.Status != 0 {
		return false
	}

	if err := s.TestServer(&server); err != nil {
		server.Status = 1
		if err := model.Update(server.ID, server); err != nil {
			fmt.Println(err.Error())
		}
		return false
	} else {
		server.Status = 0
		if err := model.Update(server.ID, server, true); err != nil {
			fmt.Println(err.Error())
		}
	}

	return true
}

func (s *ServerService) TestServer(server *model.TaskServer) error {
	if server.ConnectionType == 0 {
		switch server.Type {
		case 0:
			//密码登录
			return libs.RemoteCommandByPassword(server)
		case 1:
			//密钥登录
			return libs.RemoteCommandByKey(server)
		default:
			return errors.New("未知的登录方式")

		}
	} else if server.ConnectionType == 1 {
		if server.Type == 0 {
			//密码登录]
			return libs.RemoteCommandByTelnetPassword(server)
		} else {
			return errors.New("Telnet方式暂不支持密钥登陆！")
		}

	} else if server.ConnectionType == 2 {
		return libs.RemoteAgent(server)
	}

	return errors.New("未知错误")
}
