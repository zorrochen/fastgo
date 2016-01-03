package main
import "fmt"

type tran_module_t struct {

}


type module_fmt_t struct {
	PackName string       `xml:"packName"`
	Name     string       `xml:"name"`
	Funcs    []func_fmt_t `xml:"funcs>func"`
}

type func_fmt_t struct {
	Name         string          `xml:"name"`
	Owner        string          `xml:"owner"`
	InputParams  []inputParam_t  `xml:"inputParams>param"`
	OutputParams []outputParam_t `xml:"outputParams>param"`
	Mark         string          `xml:"mark"`
}

type inputParam_t struct {
	Name      string `xml:"name,attr"`
	ParamType string `xml:"type,attr"`
}

type outputParam_t struct {
	ParamType string `xml:"type,attr"`
}

func inputParamStr(params []inputParam_t) string {
	ret := ""
	isFirstParam := true
	for _, param := range params {
		if isFirstParam {
			isFirstParam = false
		} else {
			ret += ", "
		}
		ret += param.Name
		ret += " "
		ret += param.ParamType
	}
	return ret
}

func outputParamStr(params []outputParam_t) string {
	ret := ""
	isFirstParam := true
	for _, param := range params {
		if isFirstParam {
			isFirstParam = false
		} else {
			ret += ", "
		}
		ret += param.ParamType
	}
	return ret
}



func genModuleStr(srcFileInfo *module_fmt_t) (string, error) {
	str := ""
	str += fmt.Sprintf("package %s\n\n", srcFileInfo.PackName)

	for _, v := range srcFileInfo.Funcs {
		str += fmt.Sprintf("//%s\n", v.Mark)
		str += fmt.Sprintf("func (owner *%s_t) %s(%s) (%s) {\n\n}\n\n", v.Owner, v.Name, inputParamStr(v.InputParams), outputParamStr(v.OutputParams))
	}

	return str, nil
}
