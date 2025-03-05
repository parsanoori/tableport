package tableport

import (
	"errors"
	"fmt"
	"reflect"
)

func recursion(in interface{}, tag string, f func(key, value string, i int)) {
	// iterate over the fields of the struct and get the column names using reflect package
	t := reflect.TypeOf(in)

	// if the type is then deference it
	if t.Kind() == reflect.Ptr {
		in = reflect.ValueOf(in).Elem().Interface()
		t = reflect.TypeOf(in)
	}

	// the type should a struct
	if t.Kind() != reflect.Struct {
		return
	}

	// iterate over the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			recursion(reflect.ValueOf(in).Field(i).Interface(), tag, f)
			continue
		}
		key := t.Field(i).Tag.Get(tag)
		value := fmt.Sprintf("%v", reflect.ValueOf(in).Field(i).Interface())
		f(key, value, i)
	}
}

func flatten(in interface{}, tag string) (keys []string, values [][]string, err error) {
	// iterate over the fields of the struct and get the column names using reflect package
	t := reflect.TypeOf(in)

	// if the type is then deference it
	if t.Kind() == reflect.Ptr {
		in = reflect.ValueOf(in).Elem().Interface()
		t = reflect.TypeOf(in)
	}

	// the type should a slice or an array of structs
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return []string{}, [][]string{}, nil
	}

	// get the type of the struct
	t = t.Elem()

	// get the length of the slice or array
	l := reflect.ValueOf(in).Len()

	// if the length is 0 return error
	if l == 0 {
		err = errors.New("input should be a slice or an array of structs")
		return
	}

	// the type should be a struct
	if t.Kind() != reflect.Struct {
		err = errors.New("input should be a slice or an array of structs")
		return
	}

	// get the keys
	recursion(reflect.ValueOf(in).Index(0).Interface(), tag, func(key, value string, i int) {
		keys = append(keys, key)
	})

	// get the values
	for i := 0; i < l; i++ {
		var valuesTemp []string
		recursion(reflect.ValueOf(in).Index(i).Interface(), tag, func(key, value string, i int) {
			valuesTemp = append(valuesTemp, value)
		})
		values = append(values, valuesTemp)
	}

	return
}
