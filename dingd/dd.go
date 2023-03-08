package dingd

import (
	"errors"
	"strings"
	"sync"
	"time"
)

type DDClient struct {
	token     string
	jsToken   string
	appKey    string
	appSecret string
	agentId   int64
	userid    sync.Map
}

func (d *DDClient) SendMsg(phoneNumbers []string, msg string) error {
	var userIdList []string
	for _, phone := range phoneNumbers {
		userIdList = append(userIdList, d.getUserId(phone))
	}
	body := map[string]interface{}{
		"agent_id":    d.agentId,
		"userid_list": strings.Join(userIdList, ","),
		"msg": map[string]interface{}{
			"msgtype": "text",
			"text": map[string]interface{}{
				"content": msg,
			},
		},
	}
	_, err := d.httpPost("topapi/message/corpconversation/asyncsend_v2", body)
	return err
}

func (d *DDClient) BroadcastMsg(msg string) error {
	body := map[string]interface{}{
		"agent_id":    d.agentId,
		"to_all_user": true,
		"msg": map[string]interface{}{
			"msgtype": "text",
			"text": map[string]interface{}{
				"content": msg,
			},
		},
	}
	_, err := d.httpPost("topapi/message/corpconversation/asyncsend_v2", body)
	return err
}

func (d *DDClient) getUserId(phoneNumber string) string {
	if data, ok := d.userid.Load(phoneNumber); ok {
		return data.(string)
	}
	if data, err := d.httpPost("topapi/v2/user/getbymobile", map[string]interface{}{"mobile": phoneNumber}); err == nil {
		userId := data["result"].(map[string]interface{})["userid"].(string)
		d.userid.Store(phoneNumber, userId)
		return userId
	} else {
		return ""
	}
}

func (d *DDClient) httpGet(api string, queryParams map[string]string) {

}

func (d *DDClient) httpPost(api string, body map[string]interface{}) (map[string]interface{}, error) {
	data, err := httpPost(api, map[string]string{"access_token": d.token}, body)
	if err == nil && data["errcode"].(float64) != 0 {
		return data, errors.New(data["errmsg"].(string))
	}
	return data, err
}

func (d *DDClient) refreshToken() {
	tick := time.NewTicker(time.Hour)
	for {
		select {
		case <-tick.C:
			getToken(d.appKey, d.appSecret)
		}
	}
}

func NewClient(appKey string, appSecret string, agentId int64) *DDClient {
	if token := getToken(appKey, appSecret); token != "" {
		if jsToken := getJSToken(token); jsToken != "" {
			ret := &DDClient{
				token:     token,
				jsToken:   jsToken,
				appKey:    appKey,
				appSecret: appSecret,
				agentId:   agentId,
			}
			go ret.refreshToken()
			return ret
		}
	}

	return nil
}
