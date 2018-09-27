package handler

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func Parse(info string) GenMode {
	var m GenMode
	if strings.Contains(info, "@startuml") {
		m = HandlerGenMode(info)
	} else if strings.Contains(info, "(proxy)") {
		m = ProxyGenMode(info)
	} else {
		m = SimpleGenMode(info)
	}
	return m
}

type UmlInfo struct {
	FuncName string
	FuncNote string
}

type UmlInfoResp struct {
	mainFunc    UmlInfo
	subFuncList []UmlInfo
}

func GetUmlInfo(umldata string) *UmlInfoResp {
	ret := &UmlInfoResp{}
	r, _ := regexp.Compile("title(.+)\n")
	mstr := r.FindString(string(umldata))
	ss := strings.Split(mstr, "title")

	r, _ = regexp.Compile(`[^\(\)]+`)
	sslist := r.FindAllString(strings.TrimSpace(ss[1]), -1)
	ret.mainFunc.FuncName = strings.TrimSpace(sslist[0])
	ret.mainFunc.FuncNote = strings.TrimSpace(sslist[1])

	r, _ = regexp.Compile(`:\w+;\n.*note right:.+\n`)
	mstrlist := r.FindAllString(string(umldata), -1)
	for _, v := range mstrlist {
		ss = strings.Split(v, "note right:")
		r, _ = regexp.Compile(`\w+`)
		funname := r.FindString(ss[0])
		ret.subFuncList = append(ret.subFuncList, UmlInfo{FuncName: funname, FuncNote: strings.TrimSpace(ss[1])})
	}

	return ret
}

func SimpleGenMode(funcdata string) GenMode {
	m := SimpleMode{}
	ss := strings.Split(funcdata, "###")

	for _, v := range ss {
		onefuncinfo := strings.Split(v, "\r\n\r\n")
		var fi baseFunc
		sss := strings.Split(onefuncinfo[0], "#")
		fi.FuncName = strings.TrimSpace(sss[0])
		fi.FuncNote = strings.TrimSpace(sss[1])
		fi.FuncReqJson = onefuncinfo[1]
		fi.FuncRespJson = onefuncinfo[2]
		m.FuncList = append(m.FuncList, fi)
	}

	return &m
}

func HandlerGenMode(alldata string) GenMode {
	splited := strings.Split(alldata, "###")
	umldata := splited[0]
	funcdataList := splited[1:]

	funcMap := map[string]baseFunc{}
	for _, v := range funcdataList {
		onefuncinfo := strings.Split(v, "\r\n\r\n")
		var fi baseFunc
		fi.FuncName = strings.TrimSpace(onefuncinfo[0])
		fi.FuncReqJson = onefuncinfo[1]
		fi.FuncRespJson = onefuncinfo[2]
		funcMap[fi.FuncName] = fi
	}

	umlinfo := GetUmlInfo(umldata)
	hm := HandlerMode{}
	hm.MainFunc.FuncName = umlinfo.mainFunc.FuncName
	hm.MainFunc.FuncNote = umlinfo.mainFunc.FuncNote
	hm.MainFunc.FuncReqJson = funcMap[hm.MainFunc.FuncName].FuncReqJson
	hm.MainFunc.FuncRespJson = funcMap[hm.MainFunc.FuncName].FuncRespJson
	for _, v := range umlinfo.subFuncList {
		onefi := baseFunc{}
		onefi.FuncName = v.FuncName
		onefi.FuncNote = v.FuncNote
		onefi.FuncReqJson = funcMap[onefi.FuncName].FuncReqJson
		onefi.FuncRespJson = funcMap[onefi.FuncName].FuncRespJson
		hm.SubFuncList = append(hm.SubFuncList, onefi)
	}

	return &hm
}

func ProxyGenMode(funcdata string) GenMode {
	m := ProxyMode{}
	ss := strings.Split(funcdata, "###")

	for _, v := range ss {
		onefuncinfo := strings.Split(v, "\r\n\r\n")
		var fi ProxyFunc
		sss := strings.Split(onefuncinfo[0], "#")
		fi.FuncName = strings.TrimSpace(sss[0])
		fi.FuncNote = strings.TrimSpace(sss[1])
		fi.FuncReqMethod = strings.TrimSpace(sss[2])
		fi.FuncReqUrl = strings.TrimSpace(sss[3])
		fi.FuncReqJson = onefuncinfo[1]
		fi.FuncRespJson = onefuncinfo[2]
		m.FuncList = append(m.FuncList, fi)
	}

	return &m
}

func ToAssignExp(ptitle string, data interface{}) string {
	rst := ""
	if ptitle != "" {
		rst = fmt.Sprintf("var %s %s_t\n", ptitle, ptitle)
	} else {
		ptitle = "rst"
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		for k, v := range data.(map[string]interface{}) {
			switch reflect.TypeOf(v).Kind() {
			case reflect.String, reflect.Int, reflect.Float64, reflect.Bool:
				rst += fmt.Sprintf("%s.%s = %s\n", ptitle, FmtFieldName(k), ToTypeValue(v))
				continue
			case reflect.Slice:
				for _, vv := range v.([]interface{}) {
					switch reflect.TypeOf(vv).Kind() {
					case reflect.String:
						rst += fmt.Sprintf("%s.%s = []string{%s}\n", ptitle, FmtFieldName(k), ToGroupTypeValue(v.([]interface{})))
					case reflect.Int:
						rst += fmt.Sprintf("%s.%s = []int{%s}\n", ptitle, FmtFieldName(k), ToGroupTypeValue(v.([]interface{})))
					case reflect.Float64:
						rst += fmt.Sprintf("%s.%s = []%s{%s}\n", ptitle, FmtFieldName(k), AdaptFloat64Type(v.([]interface{})), ToGroupTypeValue(v.([]interface{})))
					case reflect.Map:
						rst += "{\n"
						rst += ToAssignExp(k, vv)
						rst += fmt.Sprintf("%s.%s = append(%s.%s, %s)\n", ptitle, FmtFieldName(k), ptitle, FmtFieldName(k), k)
						rst += "}\n"
						continue
					}
					break
				}
				continue
			default:
				rst += "\n"
				rst += ToAssignExp(k, v)
				rst += fmt.Sprintf("%s.%s = %s\n", ptitle, FmtFieldName(k), FmtFieldName(k))
				rst += "\n"
				continue
			}
		}
	}
	return rst
}
