package service

import (
	"strconv"
	"strings"

	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/utils"
)

func Menu(uid int) (map[string][]map[string]interface{}, error) {
	data := map[string][]map[string]interface{}{}

	// 左侧导航栏
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)

	if uid != 1 {
		//普通管理员
		adminAuthIds, _ := model.RoleAuthGetByIds("0")
		// adminAuthIds, _ := model.RoleAuthGetByIds(self.user.RoleIds)
		adminAuthIdArr := strings.Split(adminAuthIds, ",")
		filters = append(filters, "id__in", adminAuthIdArr)
	}

	result, _ := model.AuthGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))
	list2 := make([]map[string]interface{}, len(result))
	allow_url := ""
	i, j := 0, 0
	for _, v := range result {
		if v.AuthUrl != " " || v.AuthUrl != "/" {
			allow_url += v.AuthUrl
		}
		row := make(map[string]interface{})
		if v.Pid == 1 && v.IsShow == 1 {
			row["Id"] = int(v.Id)
			row["Sort"] = v.Sort
			row["AuthName"] = v.AuthName
			row["AuthUrl"] = utils.URI("") + v.AuthUrl
			row["Icon"] = v.Icon
			row["Pid"] = int(v.Pid)
			list[i] = row
			i++
		}

		if v.Pid != 1 && v.IsShow == 1 {
			row["Id"] = int(v.Id)
			row["Sort"] = v.Sort
			row["AuthName"] = v.AuthName
			row["AuthUrl"] = utils.URI("") + v.AuthUrl
			row["Icon"] = v.Icon
			row["Pid"] = int(v.Pid)
			list2[j] = row
			j++
		}
	}

	data["SideMenu1"] = list[:i]  //一级菜单
	data["SideMenu2"] = list2[:j] //二级菜单

	return data, nil
}

func TaskGroups(uid int, roleIDs string) (string, string) {
	// if user.RoleIds == "0" || user.Id == 1 {
	// 	return
	// }

	Filters := make([]interface{}, 0)
	Filters = append(Filters, "status", 1)

	RoleIdsArr := strings.Split(roleIDs, ",")

	RoleIds := make([]int, 0)
	for _, v := range RoleIdsArr {
		id, _ := strconv.Atoi(v)
		RoleIds = append(RoleIds, id)
	}

	Filters = append(Filters, "id__in", RoleIds)

	Result, _ := model.RoleGetList(1, 1000, Filters...)
	serverGroups := ""
	taskGroups := ""
	for _, v := range Result {
		serverGroups += v.ServerGroupIds + ","
		taskGroups += v.TaskGroupIds + ","
	}

	return strings.Trim(serverGroups, ","), strings.Trim(taskGroups, ",")
}
