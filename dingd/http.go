package dingd

import (
	"bytes"
	"encoding/json"
	"github.com/7058011439/haoqbb/Util"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	baseUrl = "https://oapi.dingtalk.com/"
)

func httpGet(api string, queryParams map[string]string) (map[string]interface{}, error) {
	params := url.Values{}
	for k, v := range queryParams {
		params.Set(k, v)
	}
	response, err := http.Get(baseUrl + api + "?" + params.Encode())
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	replyBody, _ := ioutil.ReadAll(response.Body)
	result := make(map[string]interface{})
	err = json.Unmarshal(replyBody, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func httpPost(api string, queryParams map[string]string, body map[string]interface{}) (map[string]interface{}, error) {
	requestBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", baseUrl+api, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// 设置查询参数
	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	replyBody, _ := ioutil.ReadAll(resp.Body)
	result := make(map[string]interface{})
	err = json.Unmarshal(replyBody, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getToken(appKey string, appSecret string) string {
	ret, err := httpGet("gettoken", map[string]string{"appkey": appKey, "appsecret": appSecret})
	return Util.Ternary(err == nil && ret["errcode"].(float64) == 0, ret["access_token"].(string), "").(string)
}

func getJSToken(token string) string {
	ret, err := httpGet("get_jsapi_ticket", map[string]string{"access_token": token})
	return Util.Ternary(err == nil && ret["errcode"].(float64) == 0, ret["ticket"].(string), "").(string)
}
