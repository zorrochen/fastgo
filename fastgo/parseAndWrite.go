package main
import (
	"encoding/xml"
	"io/ioutil"
	"fmt"
	"os"
	"bufio"
	"log"
)


type parse_and_write_t struct {

}

type src_file_fmt_t struct {
	PackName	string        	`xml:"packName"`
	Name	string        		`xml:"name"`
	Funcs	[]func_fmt_t        `xml:"funcs>func"`
}

type func_fmt_t struct {
	Name	string				`xml:"name"`
	Owner	string				`xml:"owner"`
	InputParams	[]inputParam_t		`xml:"inputParams>param"`
	OutputParams []outputParam_t		`xml:"outputParams>param"`
	Mark	string				`xml:"mark"`
}

type inputParam_t struct {
	Name		string	`xml:"name,attr"`
	ParamType	string	`xml:"type,attr"`
}

type outputParam_t struct {
	ParamType	string	`xml:"type,attr"`
}

func (obj *parse_and_write_t) Parse() error {
	log.Printf("parse start!")

	fileList, err := FilesInDirection(*srcFileDir)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	for _, fileName := range fileList {
		err := parse(fileName)
		if err != nil {
			log.Printf("Error:%v", err)
			return err
		}
	}

	log.Printf("parse end!")
	return nil
}

func parse(fileName string) error {
	srcFileInfo, err := readXmlFile(*srcFileDir+fileName)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	srcFileStr, err := genSrcFileStr(srcFileInfo)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	err = writeSrcFile(srcFileInfo.Name, srcFileStr)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}
	return nil
}


func readXmlFile(fileName string) (*src_file_fmt_t, error) {
	srcDat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Error:%v", err)
		return nil, err
	}

	var ret src_file_fmt_t
	err = xml.Unmarshal(srcDat, &ret)
	if err != nil {
		log.Printf("Error:%v", err)
		return nil, err
	}
	return &ret, nil
}

func writeSrcFile(fileName, srcFileStr string) error {
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

func genSrcFileStr(srcFileInfo *src_file_fmt_t) (string ,error) {
	str := ""
	str += fmt.Sprintf("package %s\n\n", srcFileInfo.PackName)

	for _, v := range srcFileInfo.Funcs {
		str += fmt.Sprintf("//%s\n", v.Mark)
		str += fmt.Sprintf("func (owner *%s_t) %s(%s) (%s) {\n\n}\n\n", v.Owner, v.Name, inputParamStr(v.InputParams), outputParamStr(v.OutputParams))
	}

	return str, nil
}

func FilesInDirection(dir string) ([]string, error) {
	fileList, err := ioutil.ReadDir(dir) //要读取的目录地址DIR，得到列表
	if err != nil {
		log.Printf("read dir error")
		return nil, err
	}

	retFileList := []string{}
	for _, file := range fileList {
		if file.IsDir() {
			continue
		}
		retFileList = append(retFileList, file.Name())
	}

	return retFileList, nil
}

//func parseInputParam(p *inputParam_t) (string, error) {
//	retStr := p.Name
//	switch p.ParamType {
//	case bool:
//		return fmt.Sprintf("%s %s")
//	}
//}
//
//func parseOnputParam(p *outputParam_t) (string, error) {
//
//}

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