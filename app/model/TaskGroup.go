/************************************************************
** @Description: model
** @Author: haodaquan
** @Date:   2018-06-10 22:24
** @Last Modified by:   haodaquan
** @Last Modified time: 2018-06-10 22:24
*************************************************************/

package model

type TaskGroup struct {
	ID          int    `xorm:"id pk" json:"id"`
	GroupName   string `xorm:"group_name" json:"group_name"`
	Description string `xorm:"description" json:"description"`
	CreatedID   int    `xorm:"create_id" json:"created_id"`
	UpdatedID   int    `xorm:"update_id" json:"updated_id"`
	CreatedAt   int64  `xorm:"create_time" json:"created_at"`
	UpdatedAt   int64  `xorm:"update_time" json:"updated_at"`
	Status      int    `xorm:"status" json:"status"`
}

func (t *TaskGroup) TableName() string {
	return "pp_task_group"
}

// func (t *TaskGroup) Update(fields ...string) error {
// 	if t.GroupName == "" {
// 		return fmt.Errorf("组名不能为空")
// 	}
// 	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func GroupAdd(obj *TaskGroup) (int64, error) {
// 	if obj.GroupName == "" {
// 		return 0, fmt.Errorf("组名不能为空")
// 	}
// 	return orm.NewOrm().Insert(obj)
// }

// func GroupGetById(id int) (*TaskGroup, error) {
// 	obj := &TaskGroup{
// 		ID: id,
// 	}
// 	err := orm.NewOrm().Read(obj)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj, nil
// }

// func GroupDelById(id int) error {
// 	_, err := orm.NewOrm().QueryTable(TableName("task_group")).Filter("id", id).Delete()
// 	return err
// }

// func GroupGetList(page, pageSize int, filters ...interface{}) ([]*TaskGroup, int64) {
// 	offset := (page - 1) * pageSize
// 	list := make([]*TaskGroup, 0)
// 	query := orm.NewOrm().QueryTable(TableName("task_group"))
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
