package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ChimeraCoder/gojson"
	"go/format"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"unicode"
)

var (
	serviceName = flag.String("srv", "", "service")
	// packageName = flag.String("pkg", "", "package")
	funcName = flag.String("func", "", "function")
	funcType = flag.Int("type", 1, "function type")
	mockFlag = flag.Bool("mock", false, "mock data switch")
)

func main() {
	flag.Parse()

	// 设置使用的CPU个数
	log.Print("start...(CPU:%d)", runtime.NumCPU())

	geninfo := &handleInfo{}
	if *funcType == 1 {
		geninfo, _ = getGenInfo(*serviceName, *funcName)
	} else if *funcType == 2 {
		geninfo, _ = getGenInfoWithUML(*serviceName, *funcName)
		if *mockFlag {
			geninfo.mainFunc.FuncMockFlag = true

			mockData := genMock(geninfo.mainFunc.FuncName, geninfo.mainFunc.FuncRespJson)
			mockData_formated_byte, _ := format.Source([]byte(mockData))
			mockData = "\n" + string(mockData_formated_byte)
			exportMock(mockData, *serviceName, *funcName, *funcType)
		}
	} else if *funcType == 3 {
		geninfo, _ = getGenInfo(*serviceName, *funcName)
	} else {
		return
	}

	genrst := ""
	if *funcType == 1 {
		genrst = "package handle\n\n"
		for _, v := range geninfo.subFuncList {
			genrst += gen(v)
		}
	} else if *funcType == 2 {
		genrst = "package handle\n\n"
		genrst += gen(geninfo.mainFunc)
		for _, v := range geninfo.subFuncList {
			genrst += gen(v)
		}
	} else if *funcType == 3 {
		genrst = "package proxy\n\n"
		importStrFmt := `import(
			"%s/proxy"
			"encoding/json"
			"errors"
			"net/http"
			"fmt"
			)`
		genrst += fmt.Sprintf(importStrFmt, *serviceName)
		genrst += "\n\n"

		for _, v := range geninfo.subFuncList {
			genrst += genProxy(v)
		}
	}

	genrst_formated_byte, _ := format.Source([]byte(genrst))
	genrst_formated := string(genrst_formated_byte)
	export(genrst_formated, *serviceName, *funcName, *funcType)

	log.Print("server stoped")
}

type handleInfo struct {
	mainFunc    funcinfo_t
	subFuncList []funcinfo_t
}

type funcinfo_t struct {
	FuncName      string
	FuncNote      string
	FuncReqJson   string
	FuncRespJson  string
	FuncMockFlag  bool
	FuncReqUrl    string
	FuncReqMethod string
}

func genMock(funcName, respJson string) string {
	genData := map[string]string{}

	genData["funcname"] = funcName
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(respJson), &data)
	if err != nil {
		fmt.Printf("json error\n")
		return ""
	}

	genData["body"] = ToAssignExp("", data)
	t, _ := template.ParseFiles("./mock.tmpl")
	b2 := &bytes.Buffer{}
	t.Execute(b2, genData)

	return b2.String()
}

func gen(fi funcinfo_t) string {
	exHeadStr := fmt.Sprintf("//================= %s =================\n", fi.FuncName)

	reqName := fi.FuncName + "Req"
	respName := fi.FuncName + "Resp"

	req := strings.NewReader(fi.FuncReqJson)
	rstreq, err := gojson.Generate(req, gojson.ParseJson, reqName, "", []string{"json"}, false, true)
	if err != nil {
		return ""
	}

	resp := strings.NewReader(fi.FuncRespJson)
	rstresp, err := gojson.Generate(resp, gojson.ParseJson, respName, "", []string{"json"}, false, true)
	if err != nil {
		fmt.Printf("err:%v", err)
		return ""
	}

	genData := map[string]string{}
	genData["funcname"] = fi.FuncName
	genData["funcnote"] = fi.FuncNote
	if fi.FuncMockFlag {
		genData["mockflag"] = "true"
	}
	t, _ := template.ParseFiles("./handle.tmpl")
	b := &bytes.Buffer{}
	t.Execute(b, genData)

	retStr := exHeadStr + string(rstreq) + "\n\n" + string(rstresp) + "\n\n" + b.String() + "\n\n"
	newStr, _ := format.Source([]byte(retStr))

	retStr = string(newStr)

	return retStr
}

func genProxy(fi funcinfo_t) string {
	exHeadStr := fmt.Sprintf("//================= %s =================\n", fi.FuncName)

	reqName := fi.FuncName + "Req"
	respName := fi.FuncName + "Resp"

	req := strings.NewReader(fi.FuncReqJson)
	rstreq, err := gojson.Generate(req, gojson.ParseJson, reqName, "", []string{"json"}, false, true)
	if err != nil {
		fmt.Printf("err:%v", err)
		return ""
	}

	resp := strings.NewReader(fi.FuncRespJson)
	rstresp, err := gojson.Generate(resp, gojson.ParseJson, respName, "", []string{"json"}, false, true)
	if err != nil {
		fmt.Printf("err:%v", err)
		return ""
	}

	genData := map[string]string{}
	genData["funcname"] = fi.FuncName
	genData["funcnote"] = fi.FuncNote
	genData["reqpath"] = fi.FuncReqUrl
	if fi.FuncReqMethod == "get" {
		genData["methodget"] = "true"
	} else if fi.FuncReqMethod == "post" {
		genData["methodpost"] = "true"
	}
	t, _ := template.ParseFiles("./proxy.tmpl")
	b := &bytes.Buffer{}
	t.Execute(b, genData)

	retStr := exHeadStr + string(rstreq) + "\n\n" + string(rstresp) + "\n\n" + b.String() + "\n\n"
	newStr, _ := format.Source([]byte(retStr))

	retStr = string(newStr)

	return retStr
}

const (
	FUNC_TYPE_HANDLE     = 1
	FUNC_TYPE_HANDLE_UML = 2
	FUNC_TYPE_PROXY      = 3
)

func export(code string, module string, FuncName string, funcType int) {
	gopath := os.Getenv("GOPATH")

	exportfile := ""
	if funcType == FUNC_TYPE_HANDLE {
		exportfile = fmt.Sprintf("%s/src/%s/handle/%s.go", gopath, module, FuncName)
		os.MkdirAll(path.Dir(exportfile), os.ModePerm)
	} else if funcType == FUNC_TYPE_PROXY {
		exportfile = fmt.Sprintf("%s/src/%s/proxy/%s/%s.go", gopath, module, FuncName, FuncName)
		os.MkdirAll(path.Dir(exportfile), os.ModePerm)
	} else {
		return
	}

	exist, _ := PathExists(exportfile)
	// 已存在，不写入
	if exist {
		fmt.Printf("allready exist.\n")
		return
	}
	writeFile(exportfile, code)
	return
}

func exportMock(code string, module string, FuncName string, funcType int) {
	gopath := os.Getenv("GOPATH")

	exportfile := ""
	if funcType == FUNC_TYPE_HANDLE {
		exportfile = fmt.Sprintf("%s/src/%s/handle/mock.go", gopath, module)
		os.MkdirAll(path.Dir(exportfile), os.ModeDir)
	} else if funcType == FUNC_TYPE_PROXY {
		exportfile = fmt.Sprintf("%s/src/%s/proxy/mock.go", gopath, module)
		os.MkdirAll(path.Dir(exportfile), os.ModeDir)
	} else {
		return
	}

	exist, _ := PathExists(exportfile)
	if !exist {
		AppendFile(exportfile, "package handle\n\n")
	}
	AppendFile(exportfile, code)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type UmlInfo struct {
	FuncName string
	FuncNote string
}

type UmlInfoResp struct {
	mainFunc    UmlInfo
	subFuncList []UmlInfo
}

func GetUmlInfo(module string, FuncName string) *UmlInfoResp {
	f := fmt.Sprintf("./gendata/%s.plantuml", FuncName)
	info, err := readFile(f)
	if err != nil {
		return nil
	}

	ret := &UmlInfoResp{}
	r, _ := regexp.Compile("title(.+)\n")
	mstr := r.FindString(string(info))
	ss := strings.Split(mstr, "title")

	r, _ = regexp.Compile(`[^\(\)]+`)
	sslist := r.FindAllString(strings.TrimSpace(ss[1]), -1)
	ret.mainFunc.FuncName = strings.TrimSpace(sslist[0])
	ret.mainFunc.FuncNote = strings.TrimSpace(sslist[1])

	r, _ = regexp.Compile(`:\w+;\n.*note right:.+\n`)
	mstrlist := r.FindAllString(string(info), -1)
	for _, v := range mstrlist {
		ss = strings.Split(v, "note right:")
		r, _ = regexp.Compile(`\w+`)
		funname := r.FindString(ss[0])
		ret.subFuncList = append(ret.subFuncList, UmlInfo{FuncName: funname, FuncNote: strings.TrimSpace(ss[1])})
	}

	return ret
}

func getGenInfoWithUML(module string, FuncName string) (*handleInfo, error) {
	f := fmt.Sprintf("./gendata/%s", FuncName)

	info, err := readFile(f)
	if err != nil {
		return nil, err
	}

	funcMap := map[string]funcinfo_t{}
	ss := strings.Split(string(info), "###")
	for _, v := range ss {
		onefuncinfo := strings.Split(v, "\n\n")
		var fi funcinfo_t
		fi.FuncName = strings.TrimSpace(onefuncinfo[0])
		fi.FuncReqJson = onefuncinfo[1]
		fi.FuncRespJson = onefuncinfo[2]
		funcMap[fi.FuncName] = fi
	}

	umlinfo := GetUmlInfo(module, FuncName)
	var hi handleInfo
	hi.mainFunc.FuncName = umlinfo.mainFunc.FuncName
	hi.mainFunc.FuncNote = umlinfo.mainFunc.FuncNote
	hi.mainFunc.FuncReqJson = funcMap[hi.mainFunc.FuncName].FuncReqJson
	hi.mainFunc.FuncRespJson = funcMap[hi.mainFunc.FuncName].FuncRespJson
	for _, v := range umlinfo.subFuncList {
		onefi := funcinfo_t{}
		onefi.FuncName = v.FuncName
		onefi.FuncNote = v.FuncNote
		onefi.FuncReqJson = funcMap[onefi.FuncName].FuncReqJson
		onefi.FuncRespJson = funcMap[onefi.FuncName].FuncRespJson
		hi.subFuncList = append(hi.subFuncList, onefi)
	}

	return &hi, nil
}

func getGenInfo(module string, FuncName string) (*handleInfo, error) {
	f := fmt.Sprintf("./gendata/%s", FuncName)

	info, err := readFile(f)
	if err != nil {
		return nil, err
	}

	var hi handleInfo
	ss := strings.Split(string(info), "###")

	for _, v := range ss {
		onefuncinfo := strings.Split(v, "\n\n")
		var fi funcinfo_t

		sss := strings.Split(onefuncinfo[0], "#")
		fi.FuncName = strings.TrimSpace(sss[0])
		fi.FuncNote = strings.TrimSpace(sss[1])
		fi.FuncReqMethod = strings.TrimSpace(sss[2])
		fi.FuncReqUrl = strings.TrimSpace(sss[3])
		fi.FuncReqJson = onefuncinfo[1]
		fi.FuncRespJson = onefuncinfo[2]
		hi.subFuncList = append(hi.subFuncList, fi)
	}

	return &hi, nil
}

func readFile(fileName string) ([]byte, error) {
	srcDat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Error:%v", err)
		return nil, err
	}

	return srcDat, nil
}

func writeFile(fileName, srcFileStr string) error {
	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(srcFileStr)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	w.Flush()
	return nil
}

func AppendFile(fileName, srcFileStr string) error {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(srcFileStr)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	w.Flush()
	return nil
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

func ToTypeValue(v interface{}) string {
	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		return fmt.Sprintf("\"%v\"", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func ToGroupTypeValue(vlist []interface{}) string {
	rstlist := []string{}
	for _, v := range vlist {
		rstlist = append(rstlist, ToTypeValue(v))
	}
	return strings.Join(rstlist, ",")
}

func AdaptFloat64Type(f64list []interface{}) string {
	for _, v := range f64list {
		if disambiguateFloatInt(v) == "float64" {
			return "float64"
		}
	}
	return "int"
}

func disambiguateFloatInt(value interface{}) string {
	const epsilon = .0001
	vfloat := value.(float64)
	if math.Abs(vfloat-math.Floor(vfloat+epsilon)) < epsilon {
		var tmp int64
		return reflect.TypeOf(tmp).Name()
	}
	return reflect.TypeOf(value).Name()
}

func FmtFieldName(s string) string {
	runes := []rune(s)
	for len(runes) > 0 && !unicode.IsLetter(runes[0]) && !unicode.IsDigit(runes[0]) {
		runes = runes[1:]
	}
	if len(runes) == 0 {
		return "_"
	}

	s = stringifyFirstChar(string(runes))
	name := lintFieldName(s)
	runes = []rune(name)
	for i, c := range runes {
		ok := unicode.IsLetter(c) || unicode.IsDigit(c)
		if i == 0 {
			ok = unicode.IsLetter(c)
		}
		if !ok {
			runes[i] = '_'
		}
	}
	s = string(runes)
	s = strings.Trim(s, "_")
	if len(s) == 0 {
		return "_"
	}
	return s
}

func lintFieldName(name string) string {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}

	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		runes := []rune(name)
		if u := strings.ToUpper(name); commonInitialisms[u] {
			copy(runes[0:], []rune(u))
		} else {
			runes[0] = unicode.ToUpper(runes[0])
		}
		return string(runes)
	}

	allUpperWithUnderscore := true
	for _, r := range name {
		if !unicode.IsUpper(r) && r != '_' {
			allUpperWithUnderscore = false
			break
		}
	}
	if allUpperWithUnderscore {
		name = strings.ToLower(name)
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word

		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))

		} else if strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}

var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"NTP":   true,
	"DB":    true,
}

var intToWordMap = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

// convert first character ints to strings
func stringifyFirstChar(str string) string {
	first := str[:1]

	i, err := strconv.ParseInt(first, 10, 8)

	if err != nil {
		return str
	}

	return intToWordMap[i] + "_" + str[1:]
}
