package SVar

import (
	"encoding/json"
)

type SVar map[string]interface{}

func (s SVar) SVar(key string) SVar {
	if ret, ok := s[key]; ok {
		return ret.(SVar)
	} else {
		ret := NewSVar()
		s[key] = ret
		return ret
	}
}

func (s SVar) String() string {
	if b, err := json.Marshal(s); err == nil {
		return string(b)
	}
	return ""
}

func NewSVar() SVar {
	var ret SVar
	ret = map[string]interface{}{}
	return ret
}

func NewSVarByData(data string) SVar {
	ret := NewSVar()
	json.Unmarshal([]byte(data), &ret)
	return ret
}
