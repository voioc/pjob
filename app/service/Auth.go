package service

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/utils"
)

type AuthService struct {
	common.Base
}

// AuthS instance
func AuthS(c *gin.Context) *AuthService {
	return &AuthService{Base: common.Base{C: c}}
}

func (s *AuthService) AuthList(page, pageSize int, filters ...interface{}) ([]*model.Auth, int64, error) {
	offset := (page - 1) * pageSize
	data := make([]*model.Auth, 0)

	// query := model.GetDB()
	// var count int
	condition := " 1 = 1 "
	if len(filters) > 0 {
		for k := 0; k < len(filters); k += 2 {
			condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
		}
	}
	// fmt.Println(condition)

	total, err := model.GetDB().Where(condition).Count(&model.Auth{})
	if err != nil {
		return nil, 0, err
	}

	if err := model.GetDB().Where(condition).Limit(pageSize, offset).Find(&data); err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (s *AuthService) AuthByID(id int) (*model.Auth, error) {
	a := new(model.Auth)

	if flag, err := model.GetDB().Where("id", id).Get(a); !flag || err != nil {
		if !flag {
			err = fmt.Errorf("auth not found")
		}

		return nil, err
	}

	return a, nil
}

//获取多个
func (s *AuthService) RoleAuthByIDs(RoleIDs string) (string, error) {
	list := make([]*model.RoleAuth, 0)
	// query := orm.NewOrm().QueryTable(TableName("uc_role_auth"))
	ids := strings.Split(RoleIDs, ",")
	// _, err = query.Filter("role_id__in", ids).All(&list, "AuthId")
	// if err != nil {
	// 	return "", err
	// }

	if err := model.GetDB().In("role_id", ids).Find(&list); err != nil {
		return "", nil
	}

	b := bytes.Buffer{}
	for _, v := range list {
		if v.AuthID != 0 && v.AuthID != 1 {
			b.WriteString(strconv.Itoa(v.AuthID))
			b.WriteString(",")
		}
	}
	AuthIDs := strings.TrimRight(b.String(), ",")
	return AuthIDs, nil
}

func (s *AuthService) Menu(uid int) (map[string][]map[string]interface{}, error) {
	data := map[string][]map[string]interface{}{}

	// 左侧导航栏
	filters := make([]interface{}, 0)
	filters = append(filters, "status = ", 1)

	if uid != 1 {
		//普通管理员
		adminAuthIds, _ := s.RoleAuthByIDs("0")
		// adminAuthIds, _ := model.RoleAuthGetByIds(self.user.RoleIds)
		adminAuthIdArr := strings.Split(adminAuthIds, ",")
		filters = append(filters, "id", adminAuthIdArr)
	}

	result, _, _ := s.AuthList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))
	list2 := make([]map[string]interface{}, len(result))
	allow_url := ""
	i, j := 0, 0
	for _, v := range result {
		if v.AuthUrl != " " || v.AuthUrl != "/" {
			allow_url += v.AuthUrl
		}
		row := make(map[string]interface{})
		if v.PID == 1 && v.IsShow == 1 {
			row["Id"] = int(v.ID)
			row["Sort"] = v.Sort
			row["AuthName"] = v.AuthName
			row["AuthUrl"] = utils.URI("") + v.AuthUrl
			row["Icon"] = v.Icon
			row["Pid"] = int(v.PID)
			list[i] = row
			i++
		}

		if v.PID != 1 && v.IsShow == 1 {
			row["Id"] = int(v.ID)
			row["Sort"] = v.Sort
			row["AuthName"] = v.AuthName
			row["AuthUrl"] = utils.URI("") + v.AuthUrl
			row["Icon"] = v.Icon
			row["Pid"] = int(v.PID)
			list2[j] = row
			j++
		}
	}

	data["SideMenu1"] = list[:i]  //一级菜单
	data["SideMenu2"] = list2[:j] //二级菜单

	return data, nil
}

func (s *AuthService) TaskGroups(uid int, roleIDs string) (string, string) {
	if uid == 1 || roleIDs == "0" {
		return "", ""
	}

	filters := make([]interface{}, 0)
	filters = append(filters, "status = ", 1)

	RoleIdsArr := strings.Split(roleIDs, ",")

	RoleIds := make([]int, 0)
	for _, v := range RoleIdsArr {
		id, _ := strconv.Atoi(v)
		RoleIds = append(RoleIds, id)
	}

	filters = append(filters, "id", RoleIds)

	result, _, _ := RoleS(s.C).RoleList(1, 1000, filters...)
	serverGroups := ""
	taskGroups := ""
	for _, v := range result {
		serverGroups += v.ServerGroupIDs + ","
		taskGroups += v.TaskGroupIDs + ","
	}

	return strings.Trim(serverGroups, ","), strings.Trim(taskGroups, ",")
}

func (s *AuthService) AuthAdd(auth *model.Auth) (int64, error) {
	return model.GetDB().Insert(auth)
}

func (s *AuthService) Update(auth *model.Auth) error {
	if _, err := model.GetDB().Where("id = ?", auth.ID).Update(auth); err != nil {
		return err
	}

	return nil
}
