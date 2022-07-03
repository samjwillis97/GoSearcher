package main

import (
	"fmt"
	"reflect"
	"strings"
)

// ToMap struct to Map[string]interface{}
func ToMap(in interface{}, tagName string) (map[string]interface{}, []string, error) {
	var fields []string
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // Non-structural return error
		return nil, nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// Traversing structure fields
	// Specify the tagName value as the key in the map; the field value as the value in the map
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			fields = append(fields, tagValue)
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, fields, nil
}

func snakeToDisplay(s string) string {
	newString := strings.Replace(s, "_", " ", -1)
	return strings.Title(newString)
}
