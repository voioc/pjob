package service

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
)

type BanService struct {
	common.Base
}

// BanS instance
func BanS(c *gin.Context) *BanService {
	return &BanService{Base: common.Base{C: c}}
}

func (s *BanService) BanList(page, pageSize int, filters ...interface{}) ([]*model.Ban, int64, error) {
	offset := (page - 1) * pageSize
	data := make([]*model.Ban, 0)

	db := model.GetDB().Where("1=1")

	in := map[string]interface{}{}
	condition := " 1=1 "
	if len(filters) > 0 {
		for k := 0; k < len(filters); k += 2 {
			// 如果是数组则单独筛出来
			if _, flag := filters[k+1].([]int); flag {
				in[filters[k].(string)] = filters[k+1]
			} else {
				condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
			}
		}
	}

	if len(in) > 0 {
		for col, v := range in {
			if col != "" {
				regex := strings.Split(col, " ")
				if len(regex) == 2 && regex[1] == "not" {
					db = db.NotIn(col, v)
				} else {
					db = db.In(col, v)
				}
			}
		}
	}

	total, err := db.Where(condition).Count(&model.Ban{})
	if err != nil {
		return nil, 0, err
	}

	if err := db.Where(condition).Limit(pageSize, offset).Find(&data); err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// 检查是否含有禁用命令
func (s *BanService) CheckCommand(command string) (string, error) {
	ban := make([]model.Ban, 0)
	if err := model.GetDB().Where("status = 0").Find(&ban); err != nil {
		return "", err
	}

	// filters := make([]interface{}, 0)
	// filters = append(filters, "status", 0)
	// ban, _ := model.BanGetList(1, 1000, filters...)

	for _, v := range ban {
		if strings.Contains(command, v.Code) {
			return v.Code, nil
		}
	}

	return "", nil
}

func (s *BanService) BanByID(id int) (*model.Ban, error) {
	data := &model.Ban{}

	if _, err := model.GetDB().Where("id = ?", id).Get(data); err != nil {
		return nil, err
	}

	if data.ID == 0 {
		return nil, fmt.Errorf("server not found")
	}

	return data, nil
}

func (s *BanService) Add(data *model.Ban) (int, error) {
	_, err := model.GetDB().Insert(data)
	return data.ID, err
}

func (s *BanService) Update(data *model.Ban, args ...bool) error {
	if len(args) > 0 && args[0] {
		if _, err := model.GetDB().Cols("status").Where("id = ?", data.ID).Update(data); err != nil {
			return err
		}
	} else {
		if _, err := model.GetDB().Where("id = ?", data.ID).Update(data); err != nil {
			return err
		}
	}

	return nil
}
