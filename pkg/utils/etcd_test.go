package utils

import (
	"fmt"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestApplyLock(t *testing.T) {

	c := make(chan os.Signal)
	signal.Notify(c)

	//insect.GlobalEtcdAddress = ""
	etcd := &EtcdClient{}

	for i := 0; i < 100; i++ {
		go func() {
			time222 := time.Now().UnixNano()
			fmt.Println(time222, etcd.ApplyLock("Test"))
		}()

	}
	operate := <-c
	fmt.Println("operate:", operate)
}
