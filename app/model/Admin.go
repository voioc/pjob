/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-16 15:42:43
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:48:17
***********************************************/
package model

import (
	"fmt"
)

type Admin struct {
	ID        int    `xorm:"id pk" json:"id"`
	LoginName string `xorm:"login_name" json:"login_name"`
	RealName  string `xorm:"real_name" json:"real_name"`
	Password  string `xorm:"password" json:"-"`
	RoleIDs   string `xorm:"role_ids" json:"role_ids"`
	Phone     string `xorm:"phone" json:"phone"`
	Email     string `xorm:"email" json:"email"`
	Dingtalk  string `xorm:"dingtalk" json:"dingtalk"`
	Wechat    string `xorm:"wechat" json:"wechat"`
	Salt      string `xorm:"salt" json:"salt"`
	LastLogin int64  `xorm:"last_login" json:"last_login"`
	LastIp    string `xorm:"last_ip" json:"last_ip"`
	Status    int    `xorm:"status" json:"status"`
	CreatedID int    `xorm:"create_id" json:"created_id"`
	UpdatedID int    `xorm:"update_id" json:"updated_id"`
	CreatedAt int64  `xorm:"create_time" json:"created_at"`
	UpdatedAt int64  `xorm:"update_time" json:"created_at"`
}

func (a *Admin) TableName() string {
	return "pp_uc_admin"
}

// func AdminAdd(a *Admin) (int64, error) {
// 	return orm.NewOrm().Insert(a)
// }

func AdminGetByName(loginName string) (*Admin, error) {
	a := new(Admin)
	// err := orm.NewOrm().QueryTable(TableName("uc_admin")).Filter("login_name", loginName).One(a)
	// if err != nil {
	// 	return nil, err
	// }
	if flag, err := GetDB().Where("login_name = ?", loginName).Get(a); !flag || err != nil {
		msg := "用户不存在"
		if err != nil {
			msg = err.Error()
		}
		return nil, fmt.Errorf(msg)
	}

	return a, nil
}

// func AdminGetList(page, pageSize int, filters ...interface{}) ([]*Admin, int64) {
// 	offset := (page - 1) * pageSize
// 	list := make([]*Admin, 0)
// 	query := orm.NewOrm().QueryTable(TableName("uc_admin"))
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

// func AdminGetById(id int) (*Admin, error) {
// 	r := new(Admin)
// 	err := orm.NewOrm().QueryTable(TableName("uc_admin")).Filter("id", id).One(r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return r, nil
// }

// func (a *Admin) Update(fields ...string) error {
// 	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func RoleAuthDelete(id int) (int64, error) {
// 	query := orm.NewOrm().QueryTable(TableName("role_auth"))
// 	return query.Filter("role_id", id).Delete()
// }

// func RoleAuthMultiAdd(ras []*RoleAuth) (n int, err error) {
// 	query := orm.NewOrm().QueryTable(TableName("role_auth"))
// 	i, _ := query.PrepareInsert()
// 	for _, ra := range ras {
// 		_, err := i.Insert(ra)
// 		if err == nil {
// 			n = n + 1
// 		}
// 	}
// 	i.Close() // 别忘记关闭 statement
// 	return n, err
// }
