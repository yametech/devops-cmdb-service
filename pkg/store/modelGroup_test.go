package store

import (
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"testing"
)

func Test_Set(t *testing.T) {
	mg := ModelGroup{Uid: "hardware"}

	//fmt.Printf("%v\n", interface{}(mg).(map[string]interface{}))
	//fmt.Println(utils.MapInterface(mg))
	fmt.Println(utils.MapInterface(mg))
}
