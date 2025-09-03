package format

import "reflect"

type ObjToMapFunc[T any] func(i T) map[string]any

func SliceObj2Map[T any](inter []T, f ObjToMapFunc[T]) ([]map[string]any, int) {
	if inter == nil {
		return nil, 0
	}
	count := 0
	ret := make([]map[string]any, 0, len(inter))
	for _, item := range inter {
		if data := f(item); data != nil {
			ret = append(ret, data)
			count++
		}
	}
	return ret, count
}

func Obj2Map[T any](inter T, f ObjToMapFunc[T]) map[string]any {
	val := reflect.ValueOf(inter)
	if !val.IsValid() {
		return nil
	}
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		if val.IsNil() {
			return nil
		}
	}
	return f(inter)
}
