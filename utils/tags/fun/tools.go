package fun

import (
	"reflect"
)

func getRefVal(callName, indexName string, indexVal any, rVal reflect.Value, rType reflect.Type) any {
	found := false
	for i := 0; i < rVal.NumField(); i++ {
		fieldVal := rVal.Field(i)
		switch fieldVal.Kind() {
		case reflect.String:
			if indexVal != nil {
				// 调用链中有根据索引在数组中查找的条件
				if !found && rType.Field(i).Name == indexName && fieldVal.Interface() == indexVal {
					found = true
				} else if found && rType.Field(i).Name == callName {
					return fieldVal.String()
				}
			} else {
				// 调用链中没有根据索引在数组中查找的条件
				if rType.Field(i).Name == callName {
					return fieldVal.String()
				}
			}
		case reflect.Int:
			if indexVal != nil {
				if !found && rType.Field(i).Name == indexName && fieldVal.Interface() == indexVal {
					found = true
				} else if found && rType.Field(i).Name == callName {
					return int(fieldVal.Int())
				}
			} else {
				if rType.Field(i).Name == callName {
					return int(fieldVal.Int())
				}
			}
		case reflect.Float32:
			if indexVal != nil {
				if !found && rType.Field(i).Name == indexName && fieldVal.Interface() == indexVal {
					found = true
				} else if found && rType.Field(i).Name == callName {
					return float32(fieldVal.Float())
				}
			} else {
				if rType.Field(i).Name == callName {
					return float32(fieldVal.Float())
				}
			}
		case reflect.Float64:
			if indexVal != nil {
				if !found && rType.Field(i).Name == indexName && fieldVal.Interface() == indexVal {
					found = true
				} else if found && rType.Field(i).Name == callName {
					return fieldVal.Float()
				}
			} else {
				if rType.Field(i).Name == callName {
					return fieldVal.Float()
				}
			}
		case reflect.Struct:
			return getRefVal(callName, indexName, indexVal, fieldVal, rType.Field(i).Type)
		case reflect.Slice:
			for j := 0; j < fieldVal.Len(); j++ {
				val := getRefVal(callName, indexName, indexVal, fieldVal.Index(j), fieldVal.Index(j).Type())
				if val != nil {
					return val
				}
			}
		}
	}
	return nil
}

func getLetters(min, max string) []string {
	start := rune(min[0])
	end := rune(max[0])
	letters := make([]string, 0)
	for i := int(start); i <= int(end); i++ {
		letters = append(letters, string(rune(i)))
	}
	return letters
}
