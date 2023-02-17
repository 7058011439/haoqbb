package String

import (
	"strings"
)

func Split(msg, sp string) []string {
	if msg == "" {
		return []string{}
	}
	if sp == "" {
		return []string{msg}
	}
	var ret []string
	arr := strings.Split(msg, string(sp[0]))

	for _, c := range arr {
		ret = append(ret, Split(c, sp[1:])...)
	}
	return ret
}

func HaveMathStr(msg string, bStr string, eStr string) (bool, string) {
	beginPos := strings.Index(msg, bStr)
	endPos := strings.Index(msg, eStr)
	if beginPos != -1 && endPos != -1 && endPos > beginPos {
		return true, msg[beginPos+len(bStr) : endPos]
	}
	return false, ""
}

func HaveBeginEnd(msg string, bStr string, eStr string) (bool, string) {
	beginPos := strings.Index(msg, bStr)
	endPos := strings.Index(msg, eStr)
	if beginPos == 0 && endPos == len(msg)-len(eStr) {
		return true, msg[beginPos+len(bStr) : endPos]
	}
	return false, ""
}

func SliceToString(space string, data ...string) string {
	ret := ""
	for index, d := range data {
		ret += d
		if index != len(data)-1 {
			ret += space
		}
	}
	return ret
}
