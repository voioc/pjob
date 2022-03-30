package service

import (
	"github.com/astaxie/beego/orm"
	"github.com/voioc/cjob/app/model"
)

func AdminGetById(id int) (*model.Admin, error) {
	r := new(model.Admin)
	err := orm.NewOrm().QueryTable(model.TableName("uc_admin")).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
