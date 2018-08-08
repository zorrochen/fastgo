# fastgo
寻找最大程度的golang工程代码自动化

## 功能
1. 依据plantUML+json，自动生成`请求结构体`、`响应结构体`、`主函数业务逻辑`、`空函数`、`注释`
2. 支持:生成第三方服务调用函数
3. 支持:生成单元测试用例（开发中...）
4. 局部更新（开发中...）
5. 文档自动生成，变更联动（开发中...）
6. db层，结合orm库，SQL自动生成（开发中...）

## 期望解决的问题
1. 节省大量编码成本
2. 保证每个函数都有注释、文档、测试用例
3. 多人协作，代码质量控制，差异控制在极小单元
4. 实现单元函数高度内聚，增强可复用性
4. 简化代码评审
5. 简化的单元测试用例，可开放测试人员管理，打破开发到测试的技术墙

## 待解决问题
1. 结构体元素的参数注释，目前没有好的方案来解决

## usage
参数说明：
```
fastgo:
  -filepath string
        元数据的文件全路径
  -srv string
        service, 默认("tmp")
```

## 举例
>客户端调用server, 获取两个包裹， 服务端分别获取两个包裹，组装后返回

**假设：filepath = ./getTwoPacks.dat, 对应的数据:**
```
@startuml
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
```

**解析规则说明:**
1. 含有“@startuml”认为是handler模式，结合uml+json生成
2. 使用“###”分割各个函数元数据，假设分割后为splitlist
3. 分割后,splitlist[0]是uml数据，其他为函数元数据
4. uml数据解析后，得到主函数和子函数相关信息

**执行:**
>./fastgo -srv test -filepath ./getTwoPacks.dat

**将在$GOPATH/src下，test项目下的export目录，自动生成代码:**
```golang
package handler

//##################################################
//                     主函数
//##################################################
//================= MyPacks =================
type MyPacksReq struct {
	MyPacksParam string `json:"myPacks_param"`
}

type MyPacksResp struct {
	Pack1 string `json:"pack1"`
	Pack2 string `json:"pack2"`
}

//获取我的包裹
func MyPacks(req MyPacksReq) (MyPacksResp, error) {
	resp := MyPacksResp{}
	innerData := innerDataMyPacks{
		req: req,
	}

	//包裹1
	innerData.reqpack1 = innerData.makepack1Req()
	resppack1, err := pack1(innerData.reqpack1)
	if err != nil {
		return MyPacksResp{}, err
	}
	innerData.resppack1 = resppack1

	//包裹2
	innerData.reqpack2 = innerData.makepack2Req()
	resppack2, err := pack2(innerData.reqpack2)
	if err != nil {
		return MyPacksResp{}, err
	}
	innerData.resppack2 = resppack2

	//组装返回数据
	resp = innerData.makeResp()

	return resp, nil
}

//##################################################
//                     子函数
//##################################################
//================= pack1 =================
type pack1Req struct {
	Pack1Param string `json:"pack1_param"`
}

type pack1Resp struct {
	Pack1 string `json:"pack1"`
}

//包裹1
func pack1(req pack1Req) (pack1Resp, error) {
	resp := pack1Resp{}
	return resp, nil
}

//================= pack2 =================
type pack2Req struct {
	Pack2Param string `json:"pack2_param"`
}

type pack2Resp struct {
	Pack2 string `json:"pack2"`
}

//包裹2
func pack2(req pack2Req) (pack2Resp, error) {
	resp := pack2Resp{}
	return resp, nil
}

//##################################################
//                  封装中间请求数据
//##################################################
//单个请求涉及的中间数据集合
type innerDataMyPacks struct {
	req MyPacksReq
	// resp MyPacksResp   //(no need)
	reqpack1  pack1Req
	resppack1 pack1Resp
	reqpack2  pack2Req
	resppack2 pack2Resp
}

//组装pack1的请求数据
func (*innerDataMyPacks) makepack1Req() pack1Req {
	return pack1Req{}
}

//组装pack2的请求数据
func (*innerDataMyPacks) makepack2Req() pack2Req {
	return pack2Req{}
}

//组装返回数据
func (*innerDataMyPacks) makeResp() MyPacksResp {
	return MyPacksResp{}
}
```
