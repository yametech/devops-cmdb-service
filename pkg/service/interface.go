package service

import (
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"sync"
)

type Service struct {
	//store.IStore
	store.Neo4jDomain
	mutex sync.Mutex
}

func (s *Service) ManualQuery(query string, properties map[string]interface{}, respObj interface{}) error {
	return s.GetSession(true).Query(query, properties, respObj)
}

func (s *Service) ManualQueryRaw(query string, properties map[string]interface{}) ([][]interface{}, error) {
	return s.GetSession(true).QueryRaw(query, properties)
}

func (s *Service) ManualExecute(query string, properties map[string]interface{}) ([][]interface{}, error) {
	return s.GetSession(false).QueryRaw(query, properties)
}
