package chat

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"testing"
)

// 用beego框架
func TestServe(t *testing.T) { // restful api 连接socket
	temp := struct {
		beego.Controller
	}{}
// 	beego.Router("api/ws",&temp,"post:函数名字") // restful api 连接socket
	ser := NewServer()
	ser.Serve(temp.Ctx.ResponseWriter, temp.Ctx.Request)
	//NewServer().RegisterHandle = func(userId uint64) {
	//	// 注册后 业务处理。。。
	//	}
}

//单元测试go test test.go -v -run="函数名字"
//基准测试go test -bench="函数名字" db_test.go -v
func TestBinder_Bind(t *testing.T) { // 连接socket后注册当前用户
	message := ` [
 	{
	"from":05,
	"to":0,
	"type":0,
	"time":123456789,
	"content":[
	{
		"text":"hello user",
		"url":"",
		"height":0,
		"width":0,
		"thumb_url":"",
		"duration":0
	}]
}
]	`
	BinaryInfo, err := json.Marshal(message)
	if err != nil {
		t.Fail()
	}
	t.Log(BinaryInfo, message) // 二进制信息发送 restful api

}
