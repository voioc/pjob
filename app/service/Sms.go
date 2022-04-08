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

type SmsAjaxReturn struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Sms struct {
	Mobiles map[string]string
	Param   map[string]string
}

var SmsChan chan *Sms
var SmsUrl string

func init() {
	SmsUrl = viper.GetString("sms.url")
	poolSize := 10 // viper.GetInt("msg.pool")

	//创建通道
	SmsChan = make(chan *Sms, poolSize)

	go func() {
		for {
			select {
			case m, ok := <-SmsChan:
				if !ok {
					return
				}
				if err := m.SendSms(); err != nil {
					logzap.Ex(context.Background(), "SendSms:", err.Error())
				}
			}
		}
	}()

}

func SendSmsToChan(mobiles map[string]string, param map[string]string) bool {
	sms := &Sms{
		Mobiles: mobiles,
		Param:   param,
	}

	select {
	case SmsChan <- sms:
		return true
	case <-time.After(time.Second * 3):
		return false
	}
}

func (s *Sms) SendSms() error {

	for _, v := range s.Mobiles {
		s.Param["mobile"] = v
		res, err := utils.HttpGet(SmsUrl, s.Param)

		if err != nil {
			log.Println(err)
			return err
		}

		ajaxData := SmsAjaxReturn{}
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
