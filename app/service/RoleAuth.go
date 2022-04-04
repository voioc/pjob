package service

import (
	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
)

type RoleAuthService struct {
	common.Base
}

// RoleAuthS instance
func RoleAuthS(c *gin.Context) *RoleAuthService {
	return &RoleAuthService{Base: common.Base{C: c}}
}

func (s *RoleAuthService) RoleAuthByID(roleID int) ([]model.RoleAuth, error) {
	data := make([]model.RoleAuth, 0)
	if err := model.GetDB().Where("role_id = ?", roleID).Find(&data); err != nil {
		// fmt.Println(err.Error())
		return nil, err
	}

	return data, nil
}

// func (s *RoleAuthService) Add(data *model.RoleAuth) (int, error) {
// 	_, err := model.GetDB().Insert(data)
// 	return int(data.RoleID), err
// }

func (s *RoleAuthService) BatchAdd(data []*model.RoleAuth) error {
	_, err := model.GetDB().Insert(data)
	return err
}

// func (s *RoleAuthService) Update(data *model.Admin, args ...bool) error {
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

func (s *RoleAuthService) Del(ids interface{}) error {
	_, flag1 := ids.([]int)
	_, flag2 := ids.([]string)

	if flag1 || flag2 {
		_, err := model.GetDB().In("role_id", ids).Delete(&model.Admin{})
		return err
	}

	return nil
}
