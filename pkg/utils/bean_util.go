package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

func Stream2Json(stream []byte) (*map[string]string, error) {
	paramMap := map[string]string{}
	err := json.Unmarshal(stream, &paramMap)
	if err != nil {
		return nil, err
	}
	return &paramMap, nil
}

func MapInterface(target interface{}) *map[string]interface{} {
	b, _ := json.Marshal(target)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return &m
}

//func ObjInterface(target map[string]interface{}) interface{} {
//	b, _ := json.Marshal(target)
//	var m interface{}
//	_ = json.Unmarshal(b, &m)
//	return m
//}

// dst：目的  src：源
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
