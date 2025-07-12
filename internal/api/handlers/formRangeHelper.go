package api

import (
	"encoding/json"
	"reflect"
	"strconv"
)

type FieldGroup struct {
	GroupLabel string
	Url        string
	Fields     []FieldConfig
}

type FieldConfig struct {
	Label string
	Name  string
	Value string
	Type  string // "text" "password" 等
}

func structToStringMap(s interface{}) []FieldConfig {
	result := []FieldConfig{}
	v := reflect.ValueOf(s)
	t := v.Type()

	// 如果传入的是指针，获取其指向的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		}

		// 获取字段值并转换为字符串
		fieldValue := v.Field(i)
		var strValue string

		switch fieldValue.Kind() {
		case reflect.String:
			strValue = fieldValue.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			strValue = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			strValue = strconv.FormatUint(fieldValue.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			strValue = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
		case reflect.Bool:
			strValue = strconv.FormatBool(fieldValue.Bool())
		default:
			// 对于不支持的类型，尝试用JSON序列化
			if jsonBytes, err := json.Marshal(fieldValue.Interface()); err == nil {
				strValue = string(jsonBytes)
			} else {
				return nil
			}
		}

		result = append(result, FieldConfig{
			Label: jsonTag,
			Name:  jsonTag,
			Value: strValue,
			Type:  "text",
		})
	}

	return result
}
