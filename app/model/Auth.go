/**********************************************
** @Des: 权限因子
** @Author: haodaquan
** @Date:   2017-09-09 20:50:36
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 21:42:08
***********************************************/
package model

type Auth struct {
	ID        int    `xorm:"id pk" json:"id"`
	AuthName  string `xorm:"auth_name" json:"role_ids"`
	AuthUrl   string `xorm:"auth_url" json:"role_ids"`
	UserID    int    `xorm:"user_id" json:"role_ids"`
	PID       int    `xorm:"pid" json:"role_ids"`
	Sort      int    `xorm:"sort" json:"sort"`
	Icon      string `xorm:"icon" json:"role_ids"`
	IsShow    int    `xorm:"is_show" json:"role_ids"`
	Status    int    `xorm:"status" json:"role_ids"`
	CreatedID int    `xorm:"create_id" json:"created_id"`
	UpdatedID int    `xorm:"update_id" json:"updated_id"`
	CreatedAt int64  `xorm:"create_time" json:"created_at"`
	UpdatedAt int64  `xorm:"update_time" json:"created_at"`
}

func (a *Auth) TableName() string {
	return "pp_uc_auth"
}

// func AuthGetList(page, pageSize int, filters ...interface{}) ([]Auth, int64) {
// 	offset := (page - 1) * pageSize
// 	list := make([]Auth, 0)
// 	query := orm.NewOrm().QueryTable(TableName("uc_auth"))
// 	if len(filters) > 0 {
// 		l := len(filters)
// 		for k := 0; k < l; k += 2 {
// 			query = query.Filter(filters[k].(string), filters[k+1])
// 		}
// 	}
// 	total, _ := query.Count()
// 	query.OrderBy("pid", "sort").Limit(pageSize, offset).All(&list)

// 	return list, total
// }

// func AuthGetListByIds(authIds string, userId int) ([]*Auth, error) {

// 	list1 := make([]*Auth, 0)
// 	var list []orm.Params
// 	//list:=[]orm.Params
// 	var err error
// 	if userId == 1 {
// 		//超级管理员
// 		_, err = orm.NewOrm().Raw("select id,auth_name,auth_url,pid,icon,is_show from pp_uc_auth where status=? order by pid asc,sort asc", 1).Values(&list)
// 	} else {
// 		_, err = orm.NewOrm().Raw("select id,auth_name,auth_url,pid,icon,is_show from pp_uc_auth where status=1 and id in("+authIds+") order by pid asc,sort asc", authIds).Values(&list)
// 	}

// 	for k, v := range list {
// 		fmt.Println(k, v)
// 	}

// 	fmt.Println(list)
// 	return list1, err
// }

// func AuthAdd(auth *Auth) (int64, error) {
// 	return orm.NewOrm().Insert(auth)
// }

// func AuthGetById(id int) (*Auth, error) {
// 	a := new(Auth)

// 	err := orm.NewOrm().QueryTable(TableName("uc_auth")).Filter("id", id).One(a)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return a, nil
// }

// func (a *Auth) Update(fields ...string) error {
// 	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
// 		return err
// 	}
// 	return nil
// }
