package core

type IObject interface {
	Set(keyPath string, value interface{}) error
	Get(keyPath string) (interface{}, error)
}
