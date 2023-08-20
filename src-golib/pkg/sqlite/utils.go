package sqlite

import (
	"reflect"
)

// CreateAnyTypeSlice interface{}转为 []interface{}
func CreateAnyTypeSlice(slice interface{}) ([]interface{}, bool) {
	val, ok := isSlice(slice)
	if !ok {
		return nil, false
	}

	sliceLen := val.Len()
	out := make([]interface{}, sliceLen)
	for i := 0; i < sliceLen; i++ {
		out[i] = val.Index(i).Interface()
	}

	return out, true
}

// isSlice 判断是否为slice数据
func isSlice(arg interface{}) (reflect.Value, bool) {
	val := reflect.ValueOf(arg)
	ok := false
	if val.Kind() == reflect.Slice {
		ok = true
	}
	return val, ok
}

func mergeMap(maps ...map[string]interface{}) map[string]interface{} {
	var merged = make(map[string]interface{}, cap(maps))
	for _, m := range maps {
		for mk, mv := range m {
			merged[mk] = mv
		}
	}
	return merged
}
