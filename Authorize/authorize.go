package Authorize

import (
	"Core/AES"
	"Core/MyEncryption"
	"Core/System"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

const (
	aesKey    = "Shun7ZiRanBa"
	fileName  = "Haoqbb.key"
	customKey = "Shehuizhu1hao"
)

type config struct {
	Guid       string
	AppName    string
	DateTime   time.Time
	RandNumber int
}

var stConfig config

func loadConfig() (ret error) {
	if fileData, err := ioutil.ReadFile(fileName); err != nil {
		ret = fmt.Errorf("加载配置文件失败, err = %v", err)
	} else {
		decryptData := MyEncryption.Decrypt(fileData, []byte(customKey))
		decrypted := AES.DecryptECB(string(decryptData), aesKey)
		if err := json.Unmarshal([]byte(decrypted), &stConfig); err != nil {
			ret = fmt.Errorf("解析配置文件失败, err = %v", err)
		}
	}
	return ret
}

func Regedit(appName string, invalid func(dateTime *time.Time, err error)) {
	machineId := System.GetMachineId()
	now := time.Now()

	if err := loadConfig(); err != nil {
		invalid(&now, fmt.Errorf("%v, 当前机器码:%v, 程序名:%v, 请联系好奇宝宝授权", err.Error(), machineId, appName))
		return
	}
	if stConfig.AppName != appName || stConfig.Guid != machineId {
		invalid(&now, fmt.Errorf("授权文件错误, 当前机器码:%v, 程序名:%v, 请联系好奇宝宝授权", machineId, appName))
		return
	} else if now.Sub(stConfig.DateTime).Seconds() > 0 {
		invalid(&now, fmt.Errorf("授权文件过期, 当前机器码:%v, 程序名:%v, 请联系好奇宝宝授权", machineId, appName))
		return
	}
	invalid(&stConfig.DateTime, nil)
	go func() {
		tick := time.NewTicker(time.Second * 10)
		for {
			<-tick.C
			now = time.Now()
			if now.Sub(stConfig.DateTime).Seconds() > 0 {
				invalid(&now, fmt.Errorf("授权文件过期, 当前机器码:%v, 程序名:%v, 请联系好奇宝宝授权", machineId, appName))
			}
		}
	}()
}

func NewKey(guid string, appName string, dateTime time.Time) (ret []byte) {
	cfg := &config{
		Guid:       guid,
		AppName:    appName,
		DateTime:   dateTime,
		RandNumber: rand.Int(),
	}
	origData, _ := json.Marshal(cfg)
	encrypted := AES.EncryptECB(string(origData), aesKey)
	ret, _ = MyEncryption.Encryption([]byte(encrypted), []byte(customKey))
	return
}
