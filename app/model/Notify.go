/************************************************************
** @Description: model
** @Author: Bee
** @Date:   2019-02-15 20:21
** @Last Modified by:   Bee
** @Last Modified time: 2019-02-15 20:21
*************************************************************/
package model

import "github.com/astaxie/beego/orm"

const (
	NotifyTplTypeSystem  = "system"
	NotifyTplTypeDefault = "default"
)

type NotifyTpl struct {
	ID        int    `xorm:"id pk" json:"id"`
	Type      string `xorm:"type" json:"type"`
	TplName   string `xorm:"tpl_name" json:"tpl_name"`
	TplType   int    `xorm:"tpl_type" json:"tpl_type"`
	Title     string `xorm:"title" json:"title"`
	Content   string `xorm:"content" json:"content"`
	Status    int    `xorm:"status" json:"status"`
	CreatedID int    `xorm:"create_id" json:"created_id"`
	UpdatedID int    `xorm:"update_id" json:"updated_id"`
	CreatedAt int64  `xorm:"create_time" json:"created_at"`
	UpdatedAt int64  `xorm:"update_time" json:"created_at"`
}

func (t *NotifyTpl) TableName() string {
	return TableName("notify_tpl")
}

// func (t *NotifyTpl) Update(fields ...string) error {
// 	if t.TplName == "" {
// 		return fmt.Errorf("模板名称不能为空")
// 	}

// 	if t.Content == "" {
// 		return fmt.Errorf("模板内容不能为空")
// 	}

// 	if t.CreatedAt == 0 {
// 		t.CreatedAt = time.Now().Unix()
// 	}

// 	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func NotifyTplAdd(obj *NotifyTpl) (int64, error) {
// 	if obj.TplName == "" {
// 		return 0, fmt.Errorf("模板名称不能为空")
// 	}
// 	if obj.Content == "" {
// 		return 0, fmt.Errorf("模板内容不能为空")
// 	}
// 	if obj.CreatedAt == 0 {
// 		obj.CreatedAt = time.Now().Unix()
// 	}
// 	return orm.NewOrm().Insert(obj)
// }

func NotifyTplGetByTplType(tpl_type int, typestr string) (NotifyTpl, error) {
	var obj NotifyTpl
	err := orm.NewOrm().QueryTable(TableName("notify_tpl")).Filter("type", typestr).Filter("tpl_type", tpl_type).Filter("status", 1).Limit(1).One(&obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

// func NotifyTplGetById(id int) (*NotifyTpl, error) {
// 	obj := &NotifyTpl{
// 		ID: id,
// 	}
// 	err := orm.NewOrm().Read(obj)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj, nil
// }

// func NotifyTplGetByTplTypeList(tpl_type int) ([]*NotifyTpl, int64, error) {
// 	list := make([]*NotifyTpl, 0)
// 	total, err := orm.NewOrm().QueryTable(TableName("notify_tpl")).Filter("tpl_type", tpl_type).Filter("status", 1).All(&list)
// 	return list, total, err
// }

// func NotifyTplDelById(id int) error {
// 	_, err := orm.NewOrm().QueryTable(TableName("notify_tpl")).Filter("id", id).Delete()
// 	return err
// }

// func NotifyTplGetList(page, pageSize int, filters ...interface{}) ([]*NotifyTpl, int64) {

// 	offset := (page - 1) * pageSize
// 	list := make([]*NotifyTpl, 0)
// 	query := orm.NewOrm().QueryTable(TableName("notify_tpl"))
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
