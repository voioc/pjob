/************************************************************
** @Description: server
** @Author: george hao
** @Date:   2018-11-29 11:13
** @Last Modified by:  george hao
** @Last Modified time: 2018-11-29 11:13
*************************************************************/
package server

import (
	"github.com/go-ini/ini"
	"github.com/voioc/cjob/agent/common"
)

var C = new(common.Conf)
var ConfPath string

func InitConfig(path string) error {

	Cfg, err := ini.Load(path)
	if err != nil {
		return err
	}

	ConfPath = path
	err = Cfg.MapTo(C)
	return err
}

func SaveConfig(key string, value string) error {
	Cfg, err := ini.Load(ConfPath)
	if err != nil {
		return err
	}
	Cfg.Section("").Key(key).SetValue(value)
	Cfg.SaveTo(ConfPath)
	InitConfig(ConfPath)
	return nil
}
