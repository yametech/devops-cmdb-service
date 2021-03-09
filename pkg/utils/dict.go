package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	reflect "reflect"
	"strings"
)

func MapInterface(target interface{}) map[string]interface{} {
	b, _ := json.Marshal(target)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}

func ObjInterface(target map[string]interface{}) interface{} {
	b, _ := json.Marshal(target)
	var m interface{}
	_ = json.Unmarshal(b, &m)
	return m
}

func SimpleConvert(dst, src interface{}) {
	byteRrc, _ := json.Marshal(src)
	err := json.Unmarshal(byteRrc, dst)
	if err != nil {
		panic(err)
	}
}

func SimpleCopyProperties(dst, src interface{}) (err error) {
	// 防止意外panic
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dst type should be a struct pointer")
	}

	// src必须为结构体或者结构体指针，.Elem()类似于*ptr的操作返回指针指向的地址反射类型
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效，说明src没有这个属性 || 属性同名但类型不同
		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}

	return nil
}

// "path":"a.b.c"
// data = {"a":{"b":{"c":123}}}
// Set(data,"a.b.c",123)
func Set(data map[string]interface{}, path string, value interface{}) {
	head, remain := shift(path)
	_, exist := data[head]
	if !exist {
		data[head] = make(map[string]interface{})
	}
	if remain == "" {
		data[head] = value
		return
	}
	Set(data[head].(map[string]interface{}), remain, value)
}

// data = {"a":{"b":{"c":123}}}
// Get(data,"a.b.c") = 123
func Get(data map[string]interface{}, path string) (value interface{}) {
	head, remain := shift(path)
	_, exist := data[head]
	if exist {
		if remain == "" {
			return data[head]
		}
		switch data[head].(type) {
		case map[string]interface{}:
			return Get(data[head].(map[string]interface{}), remain)
		}
	}
	return nil
}

// data = {"a":{"b":{"c":123}}}
// Delete(data,"a.b.c") = {"a":{"b":""}}
func Delete(data map[string]interface{}, path string) {
	head, remain := shift(path)
	_, exist := data[head]
	if exist {
		if remain == "" {
			delete(data, head)
			return
		}
		switch data[head].(type) {
		case map[string]interface{}:
			Delete(data[head].(map[string]interface{}), remain)
			return
		}
	}
	return
}

func shift(path string) (head string, remain string) {
	slice := strings.Split(path, ".")
	if len(slice) < 1 {
		return "", ""
	}
	if len(slice) < 2 {
		remain = ""
		head = slice[0]
		return
	}
	return slice[0], strings.Join(slice[1:], ".")
}
