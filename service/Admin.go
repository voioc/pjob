package service

import (
	"github.com/astaxie/beego/orm"
	"github.com/voioc/cjob/models"
)

func AdminGetById(id int) (*models.Admin, error) {
	r := new(models.Admin)
	err := orm.NewOrm().QueryTable(models.TableName("uc_admin")).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
