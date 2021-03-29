package store

import (
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
)

type ModelGroup struct {
	gogm.BaseNode
	Uid    string   `json:"uid" gogm:"unique;name=uid"`
	Name   string   `json:"name" gogm:"name=name"`
	Models []*Model `json:"model" gogm:"direction=incoming;relationship=GroupBy"`
	CommonObj
}

func (g *ModelGroup) AddModel(model *Model) {
	for _, m := range g.Models {
		if m.Uid == model.Uid {
			return
		}
	}
	g.Models = append(g.Models, model)
}
