/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-15 11:44:13
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:49:13
***********************************************/
package model

type RoleAuth struct {
	AuthID int   `xorm:"auth_id pk" orm:"pk" json:"auth_id"`
	RoleID int64 `xorm:"role_id" json:"role_id"`
}

func (ra *RoleAuth) TableName() string {
	return "pp_uc_role_auth"
}

// func RoleAuthAdd(ra *RoleAuth) (int64, error) {
// 	return orm.NewOrm().Insert(ra)
// }

// func RoleAuthBatchAdd(ras *[]RoleAuth) (int64, error) {
// 	return orm.NewOrm().InsertMulti(100, ras)
// }

// func RoleAuthGetById(id int) ([]*RoleAuth, error) {
// 	list := make([]*RoleAuth, 0)
// 	query := orm.NewOrm().QueryTable(TableName("uc_role_auth"))
// 	_, err := query.Filter("role_id", id).All(&list, "AuthId")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return list, nil
// }

// func RoleAuthDelete(id int) (int64, error) {
// 	_, err := orm.NewOrm().Raw("DELETE FROM `pp_uc_role_auth` WHERE `role_id` = ?",
// 		strconv.Itoa(id)).Exec()
// 	return 0, err
// }

// //获取多个
// func RoleAuthGetByIds(RoleIds string) (Authids string, err error) {
// 	list := make([]*RoleAuth, 0)
// 	query := orm.NewOrm().QueryTable(TableName("uc_role_auth"))
// 	ids := strings.Split(RoleIds, ",")
// 	_, err = query.Filter("role_id__in", ids).All(&list, "AuthId")
// 	if err != nil {
// 		return "", err
// 	}
// 	b := bytes.Buffer{}
// 	for _, v := range list {
// 		if v.AuthID != 0 && v.AuthID != 1 {
// 			b.WriteString(strconv.Itoa(v.AuthID))
// 			b.WriteString(",")
// 		}
// 	}
// 	Authids = strings.TrimRight(b.String(), ",")
// 	return Authids, nil
// }

// func RoleAuthMultiAdd(ras []*RoleAuth) (n int, err error) {
// 	query := orm.NewOrm().QueryTable(TableName("uc_role_auth"))
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
