package utils

import "reflect"

// ConvertMapToStruct convert MapToStruct
func ConvertMapToStruct(m map[string]interface{}, s interface{}) error {
	structValue := reflect.ValueOf(s).Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)
		fieldName := fieldType.Name

		if val, ok := m[fieldName]; ok {
			field.Set(reflect.ValueOf(val))
		}
	}

	/*
		stValue := reflect.ValueOf(s).Elem()
		sType := stValue.Type()
		for i := 0; i < sType.NumField(); i++ {
			field := sType.Field(i)
			if value, ok := m[field.Name]; ok {
				stValue.Field(i).Set(reflect.ValueOf(value))
			}
		}
	*/

	return nil
}
