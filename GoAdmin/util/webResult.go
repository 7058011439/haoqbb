package util

type WebResult struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}
