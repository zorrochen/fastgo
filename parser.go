package main

import (
	"bytes"
	"fmt"
	"github.com/ChimeraCoder/gojson"
	"go/format"
	"strings"
	"text/template"
)

type GenMode interface {
	// Parse() error
	Gen() (string, error)
	// Export() error
}

type baseFunc struct {
	// FuncOwner    string
	FuncName     string
	FuncNote     string
	FuncReqJson  string
	FuncRespJson string
	Body         string
}

func GenOneFunc(fi baseFunc) string {
	exHeadStr := fmt.Sprintf("//================= %s =================\n", fi.FuncName)

	reqName := fi.FuncName + "Req"
	respName := fi.FuncName + "Resp"

	req := strings.NewReader(fi.FuncReqJson)
	rstreq, err := gojson.Generate(req, gojson.ParseJson, reqName, "", []string{"json"}, false, true)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return ""
	}

	resp := strings.NewReader(fi.FuncRespJson)
	rstresp, err := gojson.Generate(resp, gojson.ParseJson, respName, "", []string{"json"}, false, true)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return ""
	}

	genData := map[string]string{}
	genData["funcname"] = fi.FuncName
	genData["funcnote"] = fi.FuncNote
	genData["body"] = fi.Body
	t, err := template.ParseFiles("./func.tmpl")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return ""
	}
	b := &bytes.Buffer{}
	t.Execute(b, genData)

	retStr := exHeadStr + string(rstreq) + "\n\n" + string(rstresp) + "\n\n" + b.String() + "\n\n"
	newStr, err := format.Source([]byte(retStr))
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return ""
	}

	retStr = string(newStr)
	return retStr
}

func genInvokedCode(fi baseFunc) string {
	genData := map[string]string{}
	genData["funcname"] = fi.FuncName
	genData["funcnote"] = fi.FuncNote
	t, err := template.ParseFiles("./handleBody.tmpl")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return ""
	}
	b := &bytes.Buffer{}
	t.Execute(b, genData)
	return b.String()
}

func genInnerDataCode(mainFuncName string, subFuncList []string) string {
	genData := map[string][]string{}
	genData["mainFunc"] = []string{mainFuncName}
	genData["subFuncList"] = subFuncList
	t, err := template.ParseFiles("./innerDataInit.tmpl")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return ""
	}
	b := &bytes.Buffer{}
	t.Execute(b, genData)
	return b.String()
}

func genInnerDataDefineCode(mainFuncName string) string {
	genData := map[string]string{}
	genData["mainFunc"] = mainFuncName
	t, err := template.ParseFiles("./innerDataDefine.tmpl")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return ""
	}
	b := &bytes.Buffer{}
	t.Execute(b, genData)
	return b.String()
}

func genReqMaker(mainFunc string, subFunc string) string {
	genData := map[string]string{}
	genData["mainFunc"] = mainFunc
	genData["funcname"] = subFunc
	t, err := template.ParseFiles("./reqmaker.tmpl")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return ""
	}
	b := &bytes.Buffer{}
	t.Execute(b, genData)
	return b.String()
}

func genMakeResp(mainFunc string) string {
	genData := map[string]string{}
	genData["mainFunc"] = mainFunc
	t, err := template.ParseFiles("./makeResp.tmpl")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return ""
	}
	b := &bytes.Buffer{}
	t.Execute(b, genData)
	return b.String()
}

type SimpleMode struct {
	FuncList []baseFunc
}

func (m *SimpleMode) Gen() (string, error) {
	rst := ""
	for _, v := range m.FuncList {
		rst += GenOneFunc(v)
	}
	return rst, nil
}

type HandlerMode struct {
	MainFunc    baseFunc
	SubFuncList []baseFunc
	//MockFlag  bool
}

func (m *HandlerMode) Gen() (string, error) {
	genData := map[string][]string{}

	genData["innerDataDeclareCode"] = []string{m.genInnerDataStruct()}
	m.MainFunc.Body = m.genBody()
	genData["mainFuncCode"] = []string{GenOneFunc(m.MainFunc)}

	subFuncCodeList := []string{}
	for _, v := range m.SubFuncList {
		subFuncCodeList = append(subFuncCodeList, GenOneFunc(v))
	}
	genData["subFuncCodeList"] = subFuncCodeList

	reqMakerCodeList := []string{}
	for _, v := range m.SubFuncList {
		reqMakerCodeList = append(reqMakerCodeList, genReqMaker(m.MainFunc.FuncName, v.FuncName))
	}
	genData["reqMakerCodeList"] = reqMakerCodeList
	genData["makeResponse"] = []string{genMakeResp(m.MainFunc.FuncName)}

	t, err := template.ParseFiles("./handleFormat.tmpl")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return "", err
	}
	b := &bytes.Buffer{}
	t.Execute(b, genData)
	return b.String(), nil
}

func (m *HandlerMode) genInnerDataStruct() string {
	rst := "\n"
	mainFuncName := m.MainFunc.FuncName
	subFuncList := []string{}
	for _, v := range m.SubFuncList {
		subFuncList = append(subFuncList, v.FuncName)
	}
	rst += genInnerDataCode(mainFuncName, subFuncList)
	return rst
}

func (m *HandlerMode) genBody() string {
	rst := ""
	rst += genInnerDataDefineCode(m.MainFunc.FuncName)
	rst += "\n\n"
	for _, v := range m.SubFuncList {
		rst += genInvokedCode(v)
		rst += "\n\n"
	}
	rst += "//组装返回数据\n"
	rst += "resp = innerData.makeResp()\n"
	return rst
}

type ProxyFunc struct {
	baseFunc
	FuncReqUrl    string
	FuncReqMethod string
}

type ProxyMode struct {
	FuncList []ProxyFunc
}

func (m *ProxyMode) Gen() (string, error) {
	rst := ""
	for _, v := range m.FuncList {
		v.Body = v.genBody()
		rst += GenOneFunc(v.baseFunc)
	}
	return rst, nil
}

func (fi *ProxyFunc) genBody() string {
	rst := ""

	genData := map[string]string{}
	genData["reqpath"] = fi.FuncReqUrl
	if fi.FuncReqMethod == "get" {
		genData["methodget"] = "true"
	} else if fi.FuncReqMethod == "post" {
		genData["methodpost"] = "true"
	}
	t, _ := template.ParseFiles("./proxy.tmpl")
	b := &bytes.Buffer{}
	t.Execute(b, genData)

	rst += b.String()
	return rst
}
