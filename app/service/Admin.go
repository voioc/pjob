package service

import (
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

// func (s *AdminService) AdminGetByID(id int) (*model.Admin, error) {
// 	user := new(model.Admin)
// 	if _, err := model.GetDB().Where("id = ?", id).Get(user); err != nil {
// 		// fmt.Println(err.Error())
// 		return nil, err
// 	}

// 	return user, nil
// }

// func (s *AdminService) AdminList(page, pageSize int, filters ...interface{}) ([]*model.Admin, int64, error) {
// 	offset := (page - 1) * pageSize
// 	data := make([]*model.Admin, 0)

// 	// query := model.GetDB()
// 	// var count int
// 	condition := " 1 = 1 "
// 	if len(filters) > 0 {
// 		for k := 0; k < len(filters); k += 2 {
// 			condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
// 		}
// 	}

// 	total, err := model.GetDB().Where(condition).Count(&model.Admin{})
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	if err := model.GetDB().Where(condition).Limit(pageSize, offset).OrderBy("id desc").Find(&data); err != nil {
// 		return nil, 0, err
// 	}

// 	// query := orm.NewOrm().QueryTable(TableName("task"))
// 	// if len(filters) > 0 {
// 	// 	l := len(filters)
// 	// 	for k := 0; k < l; k += 2 {
// 	// 		query = query.Filter(filters[k].(string), filters[k+1])
// 	// 	}
// 	// }

// 	return data, total, nil
// }

func (s *AdminService) AdminInfo(ids []int) ([]*model.Admin, error) {
	data := make([]*model.Admin, 0)

	db := model.GetDB().Select("id, email, phone, real_name").Where("status = 1")
	if len(ids) > 0 {
		db.In("id", ids)
	}

	if err := db.Find(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// func (s *AdminService) Add(data *model.Admin) (int, error) {
// 	_, err := model.GetDB().Insert(data)
// 	return data.ID, err
// }

// func (s *AdminService) Update(data *model.Admin, args ...bool) error {
// 	if len(args) > 0 && args[0] {
// 		if _, err := model.GetDB().Cols("status").Where("id = ?", data.ID).Update(data); err != nil {
// 			return err
// 		}
// 	} else {
// 		if _, err := model.GetDB().Where("id = ?", data.ID).Update(data); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *AdminService) Del(ids interface{}) error {
// 	_, flag1 := ids.([]int)
// 	_, flag2 := ids.([]string)

// 	if flag1 || flag2 {
// 		_, err := model.GetDB().In("id", ids).Delete(&model.Admin{})
// 		return err
// 	}

// 	return nil
// }
