package utils

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	insect "github.com/yametech/go_insect"
	"log"
	"sync"
	"time"
)

type EtcdClient struct {
}

var etcdMutex sync.Mutex
var etcd *clientv3.Client

func getDefaultClient() *clientv3.Client {
	etcdMutex.Lock()
	defer etcdMutex.Unlock()

	if etcd == nil {
		log.Println("init etcd client")
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{insect.GlobalEtcdAddress},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Fatal(err)
		}
		etcd = client
	}

	return etcd
}

func (e *EtcdClient) NewMutex(key string) *concurrency.Mutex {
	session, err := concurrency.NewSession(getDefaultClient())
	if err != nil {
		log.Fatal(err)
	}
	return concurrency.NewMutex(session, "cmdb-mutex_"+key)
}

func (e *EtcdClient) ApplyLock(key string) bool {
	key = "cmdb-lock_" + key
	lease := clientv3.NewLease(getDefaultClient())
	leaseResp, err := lease.Grant(context.TODO(), 60*10)
	if err != nil {
		log.Println(err)
		return false
	}
	leaseID := leaseResp.ID

	txn := clientv3.NewKV(getDefaultClient()).Txn(context.TODO())
	cmp := clientv3.Compare(clientv3.CreateRevision(key), "=", 0)
	put := clientv3.OpPut(key, "", clientv3.WithLease(leaseID))
	txnResp, err := txn.If(cmp).Then(put).Else().Commit()
	if err != nil {
		log.Println(err)
		return false
	}
	return txnResp.Succeeded
}

func (e *EtcdClient) DeleteLease(key string) {
	getDefaultClient().Delete(context.TODO(), "cmdb-lock_"+key)
}
