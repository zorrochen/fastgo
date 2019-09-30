package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"html"
	"strings"
	"text/template"

	"github.com/zorrochen/gojson"
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
		if IsNoneParam(v.FuncReqJson) {
			continue
		}
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

	apidoc, _ := m.GenApi()
	fmt.Println(apidoc)
	return b.String(), nil
}

func (m *HandlerMode) GenApi() (string, error) {
	// api
	// req1, reqSummary1 := jsonLD.Unmarshal(m.MainFunc.FuncReqJson)
	// resp1, respSummary1 := jsonLD.Unmarshal(m.MainFunc.FuncRespJson)
	am := ApiMeta{
		Title:       m.MainFunc.FuncName,
		Method:      "post",
		Path:        "/aaa",
		Req:         m.MainFunc.FuncReqJson,
		Resp:        m.MainFunc.FuncRespJson,
		SummaryReq:  m.MainFunc.ReqSummary,
		SummaryResp: m.MainFunc.RespSummary,
	}
	apiRst := GenApi(am)
	return apiRst, nil
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
	ReqSummary   map[string]string
	RespSummary  map[string]string
}

func IsNoneParam(data string) bool {
	var dataDecoded map[string]interface{}
	json.Unmarshal([]byte(data), &dataDecoded)

	isNone := false
	if len(dataDecoded) == 0 {
		isNone = true
	}
	return isNone
}

func GenOneFunc(fi baseFunc) string {
	exHeadStr := fmt.Sprintf("//================= %s =================\n", fi.FuncName)

	reqName := fi.FuncName + "Req"
	respName := fi.FuncName + "Resp"

	isNoneReq := IsNoneParam(fi.FuncReqJson)
	isNoneResp := IsNoneParam(fi.FuncRespJson)

	rstReq := ""
	rstResp := ""
	var err error

	if !isNoneReq {
		req := strings.NewReader(fi.FuncReqJson)
		rstreq, err := gojson.Generate(req, gojson.ParseJson, reqName, "", []string{"json"}, false, true)
		if err != nil {
			fmt.Printf("err:%v\n", err)
			return ""
		}
		rstReq = string(rstreq)
	}

	if !isNoneResp {
		resp := strings.NewReader(fi.FuncRespJson)
		rstresp, err := gojson.Generate(resp, gojson.ParseJson, respName, "", []string{"json"}, false, true)
		if err != nil {
			fmt.Printf("err:%v\n", err)
			return ""
		}
		rstResp = string(rstresp)
	}

	tmpSelect := TEMP_FUNC
	switch {
	case isNoneReq && !isNoneResp:
		tmpSelect = TEMP_FUNC_NOINPUT
	case !isNoneReq && isNoneResp:
		tmpSelect = TEMP_FUNC_NOOUTPUT
	case isNoneReq && isNoneResp:
		tmpSelect = TEMP_FUNC_NOBOTH
	}

	genData := map[string]string{}
	genData["funcname"] = fi.FuncName
	genData["funcnote"] = fi.FuncNote
	genData["body"] = fi.Body
	t, err := template.New("").Parse(tmpSelect)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return ""
	}
	b := &bytes.Buffer{}
	t.Execute(b, genData)

	retStr := exHeadStr + rstReq + "\n\n" + rstResp + "\n\n" + b.String() + "\n\n"
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

	isNoneReq := IsNoneParam(fi.FuncReqJson)
	isNoneResp := IsNoneParam(fi.FuncRespJson)
	tmpSelect := TEMP_HANDLER_BODY
	switch {
	case isNoneReq && !isNoneResp:
		tmpSelect = TEMP_HANDLER_BODY_NO_REQ
	case !isNoneReq && isNoneResp:
		tmpSelect = TEMP_HANDLER_BODY_NO_RESP
	case isNoneReq && isNoneResp:
		tmpSelect = TEMP_HANDLER_BODY_NO_BOTH
	}

	t, err := template.New("").Parse(tmpSelect)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return ""
	}
	b := &bytes.Buffer{}
	t.Execute(b, genData)
	return b.String()
}

const TEMP_HANDLER_INNER_DATA_INIT = `//单个请求涉及的中间数据集合
            type innerData{{- range .mainFunc}}{{.}}{{- end}} struct {
              {{- range .mainFunc}}
              req {{.}}Req
              // resp {{.}}Resp   //(no need)
              {{- end}}
              {{- range .subFuncList}}
              req{{.}} {{.}}Req
              resp{{.}} {{.}}Resp
              {{- end}}
            }`

// func genInnerDataCode(mainFuncName string, subFuncList []string) string {
// 	genData := map[string][]string{}
// 	genData["mainFunc"] = []string{mainFuncName}
// 	genData["subFuncList"] = subFuncList
// 	t, err := template.New("").Parse(TEMP_HANDLER_INNER_DATA_INIT)
// 	if err != nil {
// 		fmt.Printf("err:%v\n", err)
// 		return ""
// 	}
// 	b := &bytes.Buffer{}
// 	t.Execute(b, genData)
// 	return b.String()
// }

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

	subGenStr := ""
	for _, v := range m.SubFuncList {
		if !IsNoneParam(v.FuncReqJson) {
			subGenStr += fmt.Sprintf(`req%s %sReq
			`, v.FuncName, v.FuncName)
		} else {
			subGenStr += fmt.Sprintf(`//req%s %sReq  //(empty)
			`, v.FuncName, v.FuncName)
		}

		if !IsNoneParam(v.FuncRespJson) {
			subGenStr += fmt.Sprintf(`resp%s %sResp
			`, v.FuncName, v.FuncName)
		} else {
			subGenStr += fmt.Sprintf(`//resp%s %sResp  //(empty)
			`, v.FuncName, v.FuncName)
		}
	}

	genStr := fmt.Sprintf(`
	//单个请求涉及的中间数据集合
	type innerData%s struct {
	  req %sReq
	  // resp %sResp  //(no need)
	  %s
	}
	`, mainFuncName, mainFuncName, mainFuncName, subGenStr)

	fmt.Printf("%s\n", genStr)

	rst += genStr
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

type ApiMeta struct {
	Title       string
	Method      string
	Path        string
	Req         string
	Resp        string
	SummaryReq  map[string]string
	SummaryResp map[string]string
}

func JsonToIndent(s string) string {
	tempMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(s), &tempMap)
	rst, _ := json.MarshalIndent(tempMap, "    ", "")
	return string(rst)
}

func GenApi(am ApiMeta) string {
	t, _ := template.New("").Parse(TEMP_API)
	b := &bytes.Buffer{}
	t.Execute(b, &am)
	rstStr := strings.Replace(b.String(), "@1", "```", -1)
	rstUnescape := html.UnescapeString(rstStr)
	return rstUnescape
}
