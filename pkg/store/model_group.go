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

func (mg *ModelGroup) Get(uid string) error {
	query := fmt.Sprintf("match (a:ModelGroup) where a.uid = $uid return a")
	properties := map[string]interface{}{
		"uid": uid,
	}
	return GetSession(false).Query(query, properties, mg)
}

func (mg *ModelGroup) Delete(uid string) error {
	query := fmt.Sprintf("match (a:ModelGroup) where a.uid = $uid return a")
	properties := map[string]interface{}{
		"uid": uid,
	}
	session := GetSession(false)
	if err := session.Query(query, properties, mg); err != nil {
		return err
	}
	if err := session.Delete(mg); err != nil {
		return err
	}
	return nil
}
