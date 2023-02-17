package DataBase

import (
	"fmt"
	"strings"
)

func conditionToString(condition map[string]interface{}) string {
	var ret strings.Builder
	firstCondition := true
	if len(condition) != 0 {
		for field, value := range condition {
			if !firstCondition {
				ret.WriteString(" and ")
			} else {
				firstCondition = false
			}
			if field != "" {
				ret.WriteString(fmt.Sprintf("%s = \"%v\"", field, value))
			} else {
				ret.WriteString(value.(string))
			}
		}
	} else {
		ret.WriteString(fmt.Sprintf("1 = 1"))
	}

	return ret.String()
}

func insertToString(data map[string]interface{}) (string, string) {
	var fields strings.Builder
	var values strings.Builder
	firstData := true

	for field, value := range data {
		if !firstData {
			fields.WriteByte(',')
			values.WriteByte(',')
		} else {
			firstData = false
		}
		fields.WriteString(field)
		if value == nil {
			values.WriteString("Null")
		} else {
			values.WriteString(fmt.Sprintf("\"%v\"", value))
		}
	}
	return fields.String(), values.String()
}

func dataToString(data map[string]interface{}) string {
	var dataList strings.Builder
	firstData := true

	for field, value := range data {
		if !firstData {
			dataList.WriteByte(',')
		} else {
			firstData = false
		}
		dataList.WriteString(fmt.Sprintf("%v = \"%v\"", field, value))
	}

	return dataList.String()
}

func selectToString(tableName string, condition map[string]interface{}, fields ...string) string {
	sql := "select "

	for index, field := range fields {
		if index == len(fields)-1 {
			sql += " " + field
		} else {
			sql += " " + field + ","
		}
	}
	if len(fields) == 0 {
		sql += " *"
	}
	sql += " from " + tableName + " where " + conditionToString(condition)
	return sql
}
