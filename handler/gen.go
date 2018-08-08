package handler

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

//===================== simpleMode =====================
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

//===================== handleMode =====================
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

	t, err := template.New("").Parse(TEMP_HANDLER)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return "", err
	}
	b := &bytes.Buffer{}
	t.Execute(b, genData)
	return b.String(), nil
}

//===================== proxyMode =====================
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

//================================================
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
	t, err := template.New("").Parse(TEMP_FUNC)
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

func genInvokedCode(mainFunc string, fi baseFunc) string {
	genData := map[string]string{}
	genData["mainFunc"] = mainFunc
	genData["funcname"] = fi.FuncName
	genData["funcnote"] = fi.FuncNote
	t, err := template.New("").Parse(TEMP_HANDLER_BODY)
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
	t, err := template.New("").Parse(TEMP_HANDLER_INNER_DATA_INIT)
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
	t, err := template.New("").Parse(TEMP_HANDLER_INNER_DATA_DEFINE)
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
	t, err := template.New("").Parse(TEMP_REQ_MAKER)
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
	t, err := template.New("").Parse(TEMP_HANDLER_MAKE_RESP)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return ""
	}
	b := &bytes.Buffer{}
	t.Execute(b, genData)
	return b.String()
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
		rst += genInvokedCode(m.MainFunc.FuncName, v)
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

func (fi *ProxyFunc) genBody() string {
	rst := ""

	genData := map[string]string{}
	genData["reqpath"] = fi.FuncReqUrl
	if fi.FuncReqMethod == "get" {
		genData["methodget"] = "true"
	} else if fi.FuncReqMethod == "post" {
		genData["methodpost"] = "true"
	}
	t, _ := template.New("").Parse(TEMP_PROXY)
	b := &bytes.Buffer{}
	t.Execute(b, genData)

	rst += b.String()
	return rst
}

// func genMock(funcName, respJson string) string {
// 	genData := map[string]string{}
//
// 	genData["funcname"] = funcName
// 	data := map[string]interface{}{}
// 	err := json.Unmarshal([]byte(respJson), &data)
// 	if err != nil {
// 		fmt.Printf("json error\n")
// 		return ""
// 	}
//
// 	genData["body"] = ToAssignExp("", data)
// 	t, _ := template.ParseFiles("./mock.tmpl")
// 	b2 := &bytes.Buffer{}
// 	t.Execute(b2, genData)
//
// 	return b2.String()
// }
