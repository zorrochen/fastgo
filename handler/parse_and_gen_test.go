package handler

import "testing"

func TestParseAndGen(t *testing.T) {
	type args struct {
		funcdata string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		//{name: "testSimpleMode", args: args{funcdata: teststr1}},
		{name: "testHandlerMode", args: args{funcdata: teststr2}},
		//{name: "testProxyMode", args: args{funcdata: teststr3}},
	}
	for _, tt := range tests {
		if got := ParseAndGen(tt.args.funcdata); got != tt.want {
			t.Errorf("%q. ParseAndGen() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

//test:simpleMode
const teststr1 = `pack1#包裹1

{
  "uid": "001"
}

{
	"error_code": 0,
	"error_msg": "ok",
	"data": {
		"pack1": "1"
	}
}
###
pack2#包裹2

{
  "uid": "001"
}

{
	"error_code": 0,
	"error_msg": "ok",
	"data": {
		"pack2": "2"
  }
}
`

const teststr2 = `@startuml
title MyPacks(获取我的包裹)
start
:pack1;
note right: 包裹1
:pack2;
note right: 包裹2
stop
@enduml

###
MyPacks

{
  "uid": "001"
}

{
	"error_code": 0,
	"error_msg": "ok",
	"data": {
		"pack1": "1",
  	"pack2": "2"
	}
}
###
pack1

{
  "uid": "001"
}

{
	"error_code": 0,
	"error_msg": "ok",
	"data": {
		"pack1": "1"
	}
}
###
pack2

{
  "uid": "001"
}

{
	"error_code": 0,
	"error_msg": "ok",
	"data": {
		"pack2": "2"
  }
}
`

const teststr3 = `WxGetToken#获取AccessToken#get#/cgi-bin/gettoken#(proxy)

{
   "corpid": "corpid",
   "corpsecret": "corpsecret"
}

{
	"errorcode": 0,
	"errormsg": "ok",
  "access_token": "accesstoken000001",
  "expires_in": 7200
}

###
WxCorpMsgSend#企业号消息发送#post#/cgi-bin/message/send#(proxy)

{
   "touser": "UserID1|UserID2|UserID3",
   "toparty": " PartyID1 | PartyID2 ",
   "totag": " TagID1 | TagID2 ",
   "msgtype": "text",
   "agentid": 1,
   "text": {
       "content": "Holiday Request For Pony(http://xxxxx)"
   },
   "safe":0
}

{
   "errcode": 0,
   "errmsg": "ok",
   "invaliduser": "UserID1",
   "invalidparty":"PartyID1",
   "invalidtag":"TagID1"
}
`
