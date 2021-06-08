package utils

import (
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
	"io"
	"testing"
	"time"
)

func TestPool(t *testing.T) {

	readPool, err := NewGenericPool(10, 100, 1*time.Hour, func() (io.Closer, error) {
		return gogm.NewSession(true)
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(readPool.getOrCreate())
	fmt.Println(readPool.numOpen)
}
