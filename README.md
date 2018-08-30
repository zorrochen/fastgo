# fastgo
寻找理想的golang工程代码自动化

## 功能
1. 依据plantUML+json，自动生成`请求结构体`、`响应结构体`、`主函数业务逻辑`、`空函数`、`注释`
2. 支持:生成第三方服务调用函数
3. 支持:生成单元测试用例
4. 局部更新（开发中...）
5. 文档自动生成，变更联动（开发中...）
6. db层，结合orm库，SQL自动生成（开发中...）

## 期望解决的问题
1. 节省编码成本
2. 保证每个函数都有注释、文档、测试用例
3. 多人协作，代码质量控制，差异控制在极小单元
4. 实现单元函数高度内聚，增强可复用性
4. 简化代码评审
5. 规范化的单元测试用例，可开放测试人员管理，打破开发到测试的技术墙，契合TDD的思想

## 待解决问题
1. 结构体元素的参数注释，目前没有好的方案来解决

## usage
参数说明：
```
fastgo:
  -filepath string
            元数据的文件全路径
  -t bool
            是否生成单元测试用例, 默认(false)
```
编译：
>make  

生成handler模式代码，及测试用例
>./fastgo -filepath ./example/handler -t

生成simple模式代码
>./fastgo -filepath ./example/simple  

生成proxy模式代码
>./fastgo -filepath ./example/proxy  

## 依赖
* [github.com/ChimeraCoder/gojson](http://github.com/ChimeraCoder/gojson)
* [github.com/cweill/gotests](http://github.com/cweill/gotests)
