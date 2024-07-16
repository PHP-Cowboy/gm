package mapHelper

import "reflect"

func StructToMap(obj interface{}) map[string]interface{} {
	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

	data := make(map[string]interface{})
	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		fieldName := objType.Field(i).Name
		data[fieldName] = field.Interface()
	}
	return data
}
