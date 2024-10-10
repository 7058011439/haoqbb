package Http

import (
	"encoding/json"
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"io"
	"io/ioutil"
	"net/http"
	URL "net/url"
	"strings"
)

type Head struct {
	data map[string]string
}

func (h *Head) Add(key, value string) *Head {
	if h.data == nil {
		h.init()
	}
	h.data[key] = value
	return h
}

func (h *Head) Del(key string) *Head {
	if h.data != nil {
		delete(h.data, key)
	}
	return h
}

func (h *Head) AddBatch(data map[string]string) {
	for k, v := range data {
		h.data[k] = v
	}
}

func (h *Head) Value() map[string]string {
	return h.data
}

func (h *Head) init() {
	h.data = map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	}
}

func NewHead(data map[string]string) *Head {
	ret := &Head{
		data: data,
	}
	if data == nil {
		ret.init()
	}
	return ret
}

type Body struct {
	data map[string]interface{}
}

func (b *Body) Add(key string, value interface{}) *Body {
	if b.data == nil {
		b.init()
	}
	b.data[key] = value
	return b
}

func (b *Body) Del(key string) *Body {
	if b.data != nil {
		delete(b.data, key)
	}
	return b
}

func (b *Body) Value() map[string]interface{} {
	return b.data
}

func (b *Body) init() {
	b.data = map[string]interface{}{}
}

func NewBody(data map[string]interface{}) *Body {
	ret := &Body{
		data: data,
	}
	if ret.data == nil {
		ret.init()
	}
	return ret
}

type Params struct {
	data map[string]interface{}
}

func (p *Params) Add(key string, value interface{}) *Params {
	if p.data == nil {
		p.init()
	}
	p.data[key] = value
	return p
}

func (p *Params) Del(key string) *Params {
	if p.data != nil {
		delete(p.data, key)
	}
	return p
}

func (p *Params) Value() map[string]interface{} {
	return p.data
}

func (p *Params) init() {
	p.data = map[string]interface{}{}
}

func NewParam(data map[string]interface{}) *Params {
	ret := &Params{
		data: data,
	}
	if ret.data == nil {
		ret.init()
	}
	return ret
}

var client *http.Client

func init() {
	//urlProxy, _ := URL.Parse("http://127.0.0.1:7890")
	client = &http.Client{
		//Transport: &http.Transport{
		//	Proxy:             http.ProxyURL(urlProxy),
		//	DisableKeepAlives: true,
		//},
	}
}

func GetHttpAsync(url string, param *Params, head *Head, callback func(map[string]interface{}, error, ...interface{}), backData ...interface{}) {
	go func() {
		resBytes, err := GetHttpSync(url, param, head)
		result := make(map[string]interface{})
		if resBytes != nil {
			err := json.Unmarshal(resBytes, &result)
			if err != nil {
				Log.Error(err.Error())
			}
		}
		if callback != nil {
			callback(result, err, backData...)
		}
	}()
}

func GetHttpSync(url string, param *Params, head *Head) ([]byte, error) {
	// 将 params 参数转换为 URL 查询字符串
	queryString := ""
	if param != nil {
		values := URL.Values{}
		for key, value := range param.data {
			values.Add(key, fmt.Sprintf("%v", value))
		}
		queryString = values.Encode()
	}

	// 构建完整的 URL
	if queryString != "" {
		url += "?" + queryString
	}

	req, _ := http.NewRequest("GET", url, nil)
	if head != nil {
		for k, v := range head.data {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		Log.Error("http get error %s %s", err, url)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Log.Error("http get io.ReadAll error %s", err)
		return nil, err
	}

	return body, nil
}

func PostHttpAsync(url string, head *Head, body *Body, callback func(map[string]interface{}, error, ...interface{}), backData ...interface{}) {
	go func() {
		result, err := PostHttpSync(url, head, body)
		if callback != nil {
			callback(result, err, backData...)
		}
	}()
}

func PostHttpSync(url string, head *Head, body *Body) (map[string]interface{}, error) {
	if head == nil {
		head = NewHead(nil)
	}
	if body == nil {
		body = NewBody(nil)
	}
	bodyData := URL.Values{}
	for key, value := range body.data {
		bodyData[key] = []string{fmt.Sprintf("%v", value)}
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(bodyData.Encode()))
	if err != nil {
		Log.Error("Failed to http.NewRequest on PostHttpAsync, err = %v", err)
		return nil, err
	}
	for k, v := range head.data { //解析header
		str := fmt.Sprintf("%v", v)
		req.Header.Set(k, str)
	}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error(err.Error())
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			Log.Error(err.Error())
		}
	}()
	replyBody, _ := ioutil.ReadAll(resp.Body)
	result := make(map[string]interface{})
	err = json.Unmarshal(replyBody, &result)
	if err != nil {
		Log.Error(err.Error())
		return nil, err
	}
	return result, nil
}
