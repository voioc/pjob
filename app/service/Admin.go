package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
)

type AdminService struct {
	common.Base
}

// AdminS instance
func AdminS(c *gin.Context) *AdminService {
	return &AdminService{Base: common.Base{C: c}}
}

func (s *AdminService) AdminGetByID(id int) (*model.Admin, error) {
	user := new(model.Admin)
	if _, err := model.GetDB().Where("id = ?", id).Get(user); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return user, nil
}

func (s *AdminService) AdminList(page, pageSize int, filters ...interface{}) ([]*model.Admin, int64, error) {
	offset := (page - 1) * pageSize
	data := make([]*model.Admin, 0)

	// query := model.GetDB()
	// var count int
	condition := " 1 = 1 "
	if len(filters) > 0 {
		for k := 0; k < len(filters); k += 2 {
			condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
		}
	}

	total, err := model.GetDB().Where(condition).Count(&model.Admin{})
	if err != nil {
		return nil, 0, err
	}

	if err := model.GetDB().Where(condition).Limit(pageSize, offset).OrderBy("id desc").Find(&data); err != nil {
		return nil, 0, err
	}

	// query := orm.NewOrm().QueryTable(TableName("task"))
	// if len(filters) > 0 {
	// 	l := len(filters)
	// 	for k := 0; k < l; k += 2 {
	// 		query = query.Filter(filters[k].(string), filters[k+1])
	// 	}
	// }

	return data, total, nil
}
