package store

import "github.com/yametech/devops-cmdb-service/pkg/core"

type IStore interface {
	List(string) ([]core.IObject, error)
	DeepSearch(obj core.IObject, edge string) ([]core.IObject, error)
	Get(string) (core.IObject, error)
	Put(core.IObject) error
	Update(core.IObject) error
	Delete(core.IObject) error
}
