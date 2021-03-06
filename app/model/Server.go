/************************************************************
** @Description: model
** @Author: haodaquan
** @Date:   2018-06-09 16:11
** @Last Modified by:   haodaquan
** @Last Modified time: 2018-06-09 16:11
*************************************************************/
package model

type TaskServer struct {
	ID             int    `xorm:"id pk" json:"id"`
	GroupID        int    `xorm:"group_id" json:"group_id"`
	ConnectionType int    `xorm:"connection_type" json:"connection_type"`
	ServerName     string `xorm:"server_name" json:"server_name"`
	ServerAccount  string `xorm:"server_account" json:"server_account"`
	ServerOuterIP  string `xorm:"server_outer_ip" json:"server_outer_ip"`
	ServerIP       string `xorm:"server_ip" json:"server_ip"`
	Port           int    `xorm:"port" json:"port"`
	Password       string `xorm:"password" json:"password"`
	PrivateKeySrc  string `xorm:"private_key_src" json:"private_key_src"`
	PublicKeySrc   string `xorm:"public_key_src" json:"public_key_src"`
	Type           int    `xorm:"type" json:"type"`
	Detail         string `xorm:"detail" json:"detail"`
	CreatedAt      int64  `xorm:"create_time" json:"created_at"`
	UpdatedAt      int64  `xorm:"update_time" json:"created_at"`
	Status         int    `xorm:"status" json:"status"`
}

func (t *TaskServer) TableName() string {
	return "pp_task_server"
}

// func (t *TaskServer) Update(fields ...string) error {
// 	if t.ServerName == "" {
// 		return fmt.Errorf("服务器名不能为空")
// 	}
// 	if t.ServerIP == "" {
// 		return fmt.Errorf("服务器IP不能为空")
// 	}

// 	if t.ServerAccount == "" {
// 		return fmt.Errorf("登录账户不能为空")
// 	}

// 	if t.Type == 0 && t.Password == "" {
// 		return fmt.Errorf("服务器密码不能为空")
// 	}

// 	if t.Type == 1 && t.PrivateKeySrc == "" {
// 		return fmt.Errorf("私钥不能为空")
// 	}

// 	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func TaskServerAdd(obj *TaskServer) (int64, error) {
// 	if obj.ServerName == "" {
// 		return 0, fmt.Errorf("服务器名不能为空")
// 	}
// 	if obj.ServerIP == "" {
// 		return 0, fmt.Errorf("服务器IP不能为空")
// 	}

// 	if obj.ServerAccount == "" {
// 		return 0, fmt.Errorf("登录账户不能为空")
// 	}

// 	if obj.Type == 0 && obj.Password == "" {
// 		return 0, fmt.Errorf("服务器密码不能为空")
// 	}

// 	if obj.Type == 1 && obj.PrivateKeySrc == "" {
// 		return 0, fmt.Errorf("私钥不能为空")
// 	}
// 	return orm.NewOrm().Insert(obj)
// }

// func TaskServerGetById(id int) (*TaskServer, error) {
// 	obj := &TaskServer{
// 		ID: id,
// 	}
// 	err := orm.NewOrm().Read(obj)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj, nil
// }

// func TaskServerForActuator(serverIp string, port int) int {
// 	serverFilters := make([]interface{}, 0)
// 	serverFilters = append(serverFilters, "status__in", []int{0, 1})
// 	serverFilters = append(serverFilters, "server_ip", serverIp)
// 	serverFilters = append(serverFilters, "port", port)

// 	server, _ := TaskServerGetList(1, 1, serverFilters...)

// 	if len(server) == 1 {
// 		return server[0].ID
// 	}
// 	return 0
// }

// //
// func TaskServerGetByIds(ids string) ([]*TaskServer, int64) {

// 	serverFilters := make([]interface{}, 0)
// 	//serverFilters = append(serverFilters, "status", 1)

// 	TaskServerIdsArr := strings.Split(ids, ",")
// 	TaskServerIds := make([]int, 0)
// 	for _, v := range TaskServerIdsArr {
// 		id, _ := strconv.Atoi(v)
// 		TaskServerIds = append(TaskServerIds, id)
// 	}
// 	serverFilters = append(serverFilters, "id__in", TaskServerIds)
// 	return TaskServerGetList(1, 1000, serverFilters...)
// }

// func TaskServerDelById(id int) error {
// 	_, err := orm.NewOrm().QueryTable(TableName("task_server")).Filter("id", id).Delete()
// 	return err
// }

// func TaskServerGetList(page, pageSize int, filters ...interface{}) ([]*TaskServer, int64) {

// 	offset := (page - 1) * pageSize
// 	list := make([]*TaskServer, 0)
// 	query := orm.NewOrm().QueryTable(TableName("task_server"))
// 	if len(filters) > 0 {
// 		l := len(filters)
// 		for k := 0; k < l; k += 2 {
// 			query = query.Filter(filters[k].(string), filters[k+1])
// 		}
// 	}
// 	total, _ := query.Count()
// 	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
// 	return list, total
// }
