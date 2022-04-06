/************************************************************
** @Description: model
** @Author: haodaquan
** @Date:   2018-06-10 19:51
** @Last Modified by:   haodaquan
** @Last Modified time: 2018-06-10 19:51
*************************************************************/
package model

type Ban struct {
	ID        int    `xorm:"id pk" json:"id"`
	Code      string `xorm:"code" json:"code"`
	CreatedAt int64  `xorm:"create_time" json:"created_at"`
	UpdatedAt int64  `xorm:"update_time" json:"updated_at"`
	Status    int    `xorm:"status" json:"status"`
}

func (t *Ban) TableName() string {
	return "pp_task_ban"
}

// func (t *Ban) Update(fields ...string) error {
// 	if t.Code == "" {
// 		return fmt.Errorf("命令不能为空")
// 	}
// 	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func BanAdd(obj *Ban) (int64, error) {
// 	if obj.Code == "" {
// 		return 0, fmt.Errorf("命令不能为空")
// 	}
// 	return orm.NewOrm().Insert(obj)
// }

// func BanGetById(id int) (*Ban, error) {
// 	obj := &Ban{
// 		ID: id,
// 	}
// 	err := orm.NewOrm().Read(obj)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj, nil
// }

// func BanDelById(id int) error {
// 	_, err := orm.NewOrm().QueryTable(TableName("task_ban")).Filter("id", id).Delete()
// 	return err
// }

// func BanGetList(page, pageSize int, filters ...interface{}) ([]*Ban, int64) {
// 	offset := (page - 1) * pageSize
// 	list := make([]*Ban, 0)
// 	query := orm.NewOrm().QueryTable(TableName("task_ban"))
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
