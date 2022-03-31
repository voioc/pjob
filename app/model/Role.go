/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-14 15:24:51
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:48:52
***********************************************/
package model

import (
	"github.com/astaxie/beego/orm"
)

type Role struct {
	Id             int    `xorm:"id pk" json:"id"`
	RoleName       string `xorm:"role_name" json:"role_name"`
	Detail         string `xorm:"detail" json:"detail"`
	ServerGroupIDs string `xorm:"server_group_ids" json:"server_group_ids"`
	TaskGroupIDs   string `xorm:"task_group_ids" json:"task_group_ids"`
	Status         int    `xorm:"status" json:"status"`
	CreatedID      int    `xorm:"create_id" json:"created_id"`
	UpdatedID      int    `xorm:"update_id" json:"updated_id"`
	CreatedAt      int64  `xorm:"create_time" json:"created_at"`
	UpdatedAt      int64  `xorm:"update_time" json:"created_at"`
}

func (a *Role) TableName() string {
	return TableName("uc_role")
}

func RoleGetList(page, pageSize int, filters ...interface{}) ([]*Role, int64) {
	offset := (page - 1) * pageSize
	list := make([]*Role, 0)
	query := orm.NewOrm().QueryTable(TableName("uc_role"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	return list, total
}

func RoleAdd(role *Role) (int64, error) {
	id, err := orm.NewOrm().Insert(role)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func RoleGetById(id int) (*Role, error) {
	r := new(Role)
	err := orm.NewOrm().QueryTable(TableName("uc_role")).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Role) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(r, fields...); err != nil {
		return err
	}
	return nil
}
