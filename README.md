# fastgo
寻找最大程度的golang工程代码自动化

## 功能
1. 依据plantUML+json，自动生成`请求结构体`、`响应结构体``空函数`、`注释`
2. 第三方服务调用函数自动生成（开发中...）
3. 自动生成test（开发中...）
4. 局部更新（开发中...）
5. 文档自动生成，变更联动（开发中...）
6. db层，结合orm库，SQL自动生成（开发中...）

## 优点
1. json数据写完，框架代码自动生成，且可返回mock数据
2. json数据写完，wiki, markdown文档自动生成
2. json数据写完，单元测试用例自动生成
3. 多人协作开发，都能保证主体代码规范一致
4. 节省大量编码成本
5. 统一化极简化的单元测试用例，可开放测试人员管理，打破开发到测试的技术墙

## usage

参数说明：func指函数名；mock判断是否生成mock数据；srv指定项目; type指定是处理函数还是第三方请求函数
```
go run main.go -v:
  -func string
        function
  -mock
        mock data switch
  -srv string
        service
  -type int
        function type, 1:handle 2:proxy
```

举例： 实现add(a,b)函数

gendata/testFunc.plantuml文件：
```
@startuml

title add(整数相加)

@enduml
```

gendata/testFunc文件：
```
add

{
  "a": 1,
  "b": 2
}

{
  "c": 3
}
```
执行：go run main.go -srv fastgo -func testFunc

将在$GOPATH/src下，fastgo项目的gendata目录，寻找testFunc文件作为初始数据，然后在handle目录，自动生成代码：
```
//================= add =================
type addReq struct {
	A int64 `json:"a"`
	B int64 `json:"b"`
}

type addResp struct {
	C int64 `json:"c"`
}

//整数相加
func add(req addReq) (*addResp, error) {
	rst := &addResp{}

	return rst, nil
}
```
