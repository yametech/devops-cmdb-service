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

func (mg *ModelGroup) Save(session *gogm.Session) error {
	mg.CreateTime = time.Now().Unix()
	mg.UpdateTime = time.Now().Unix()
	return session.Save(mg)
}

func (mg *ModelGroup) Update(session *gogm.Session) error {
	mg.UpdateTime = time.Now().Unix()
	return session.Save(mg)
}

func (mg *ModelGroup) Get(session *gogm.Session, uuid string) error {
	query := fmt.Sprintf("match (a:ModelGroup) where a.uuid = $uuid return a")
	properties := map[string]interface{}{
		"uuid": uuid,
	}
	return session.Query(query, properties, mg)
}

func (mg *ModelGroup) Delete(session *gogm.Session) error {
	return session.Delete(mg)
}
