package structUtil

import (
	"errors"
	"reflect"
)

// 结构体转map
func StructToMap(s interface{}) (mp map[string]interface{}, err error) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 确保传入的是一个结构体
	if val.Kind() != reflect.Struct {
		err = errors.New("StructToMap expects a struct")
		return
	}

	typ := val.Type()
	numFields := val.NumField()
	mp = make(map[string]interface{}, numFields)

	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()
		mp[field.Name] = value
	}

	return
}
