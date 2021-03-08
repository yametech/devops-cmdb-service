package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestGetModelAttributeList(t *testing.T) {
	baseUrl := "http://127.0.0.1:8080"
	// 表单数据
	//contentType := "application/x-www-form-urlencoded"
	//data := "name=枯藤&age=18"
	// json
	contentType := "application/json"
	data := `{"modelUid":"host"}`
	resp, err := http.Post(baseUrl+"/cmdb/web/resource/model-attribute-list", contentType, strings.NewReader(data))
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed,err:%v\n", err)
		return
	}
	fmt.Println(string(b))
}
