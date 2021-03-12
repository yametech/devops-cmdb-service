package store

import (
	"fmt"
	"github.com/mindstand/gogm"
	"time"
)

type ModelGroup struct {
	gogm.BaseNode
	Uid    string   `json:"uid" gogm:"unique;name=uid"`
	Name   string   `json:"name" gogm:"name=name"`
	Models []*Model `json:"model" gogm:"direction=incoming;relationship=GroupBy"`
	CommonObj
}

func (mg *ModelGroup) Save() error {
	mg.CreateTime = time.Now().Unix()
	mg.UpdateTime = time.Now().Unix()
	return GetSession(false).Save(mg)
}

func (mg *ModelGroup) Update() error {
	mg.UpdateTime = time.Now().Unix()
	return GetSession(false).Save(mg)
}

func (mg *ModelGroup) Get(uuid string) error {
	query := fmt.Sprintf("match (a:ModelGroup) where a.uuid = $uuid return a")
	properties := map[string]interface{}{
		"uuid": uuid,
	}
	return GetSession(false).Query(query, properties, mg)
}

func (mg *ModelGroup) Delete() error {
	return GetSession(false).Delete(mg)
}
