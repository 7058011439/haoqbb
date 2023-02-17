package IHttp

type IHttp interface {
	GetHttpAsync(url string, header map[string]string, fun func(getData map[string]interface{}, backData ...interface{}), backData ...interface{})
	PostHttpAsync(url string, header map[string]string, body map[string]interface{}, callback func(map[string]interface{}, ...interface{}), backData ...interface{})
	GetName() string
}

type http struct {
	i map[string]IHttp
}

var h = http{i: make(map[string]IHttp)}

func SetHttpAgent(iHttp IHttp) {
	h.i[iHttp.GetName()] = iHttp
}

func GetAsync(serviceName string, url string, header map[string]string, fun func(getData map[string]interface{}, backData ...interface{}), backData ...interface{}) {
	h.i[serviceName].GetHttpAsync(url, header, fun, backData...)
}

func PostHttpAsync(serviceName string, url string, header map[string]string, body map[string]interface{}, callback func(map[string]interface{}, ...interface{}), backData ...interface{}) {
	h.i[serviceName].PostHttpAsync(url, header, body, callback, backData...)
}
