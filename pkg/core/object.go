package core

import "github.com/yametech/devops-cmdb-service/pkg/utils"

type IObject interface {
	// Set  o = {} -> Set(o,"a.b","123") -> {"a":{"b":"123"}}
	Set(keyPath string, value interface{})
	// Get
	Get(keyPath string) interface{}
	// Delete
	Delete(keyPath string)
}

type Object = map[string]interface{}

func (o Object) Set(keyPath string, value interface{}) {
	utils.Set(o, keyPath, value)
}

func (o Object) Get(keyPath string) interface{} {
	return utils.Get(o, keyPath)
}

func (o Object) Delete(keyPath string) {
	utils.Delete(o, keyPath)
}
