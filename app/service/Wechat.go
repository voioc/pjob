package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/voioc/cjob/utils"
	"github.com/voioc/coco/logzap"
)

type WechatAjaxReturn struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Wechat struct {
	Accounts map[string]string
	Param    map[string]string
}

var WechatChan chan *Wechat
var WechatUrl string

func init() {
	WechatUrl = viper.GetString("wechat.url")
	poolSize := 10 // viper.GetInt("wechat.pool")

	//创建通道
	WechatChan = make(chan *Wechat, poolSize)

	go func() {
		for {
			select {
			case m, ok := <-WechatChan:
				if !ok {
					return
				}
				if err := m.SendWechat(); err != nil {
					logzap.Ex(context.Background(), "SendWechat:", err.Error())
				}
			}
		}
	}()

}

func SendWechatToChan(accounts map[string]string, param map[string]string) bool {
	wechat := &Wechat{
		Accounts: accounts,
		Param:    param,
	}

	select {
	case WechatChan <- wechat:
		return true
	case <-time.After(time.Second * 3):
		return false
	}
}

func (s *Wechat) SendWechat() error {

	for _, v := range s.Accounts {
		s.Param["account"] = v
		res, err := utils.HttpGet(WechatUrl, s.Param)

		if err != nil {
			log.Println(err)
			return err
		}

		ajaxData := WechatAjaxReturn{}
		jsonErr := json.Unmarshal([]byte(res), &ajaxData)

		if jsonErr != nil {
			return jsonErr
		}

		if ajaxData.Status != 200 {
			return errors.Errorf("msg %s", ajaxData.Message)
		}

		return nil

	}
	return nil
}
