package handler

import (
	"fmt"
	"os"
	"path"
)

const (
	TYPE_SIMPLE = 1
	TYPE_HANDLE = 2
	TYPE_PROXY  = 3
)

func Export(code string, module string, FuncName string) {
	gopath := os.Getenv("GOPATH")

	exportfile := fmt.Sprintf("%s/src/%s/export/%s.go", gopath, module, FuncName)
	exist, _ := PathExists(exportfile)
	// 已存在，不写入
	if exist {
		fmt.Printf("already exist.\n")
		return
	}
	os.MkdirAll(path.Dir(exportfile), os.ModePerm)
	writeFile(exportfile, code)
	return
}

// func Export(code string, module string, FuncName string, funcType int) {
// 	gopath := os.Getenv("GOPATH")
//
// 	exportfile := ""
// 	if funcType == TYPE_SIMPLE {
// 		exportfile = fmt.Sprintf("%s/src/%s/handler/%s.go", gopath, module, FuncName)
// 		os.MkdirAll(path.Dir(exportfile), os.ModePerm)
// 	} else if funcType == TYPE_HANDLE {
// 		exportfile = fmt.Sprintf("%s/src/%s/handler/%s.go", gopath, module, FuncName)
// 		os.MkdirAll(path.Dir(exportfile), os.ModePerm)
// 	} else if funcType == TYPE_PROXY {
// 		exportfile = fmt.Sprintf("%s/src/%s/proxy/%s/%s.go", gopath, module, FuncName, FuncName)
// 		os.MkdirAll(path.Dir(exportfile), os.ModePerm)
// 	} else {
// 		return
// 	}
//
// 	exist, _ := PathExists(exportfile)
// 	// 已存在，不写入
// 	if exist {
// 		fmt.Printf("already exist.\n")
// 		return
// 	}
// 	writeFile(exportfile, code)
// 	return
// }

// func exportMock(code string, module string, FuncName string, funcType int) {
// 	gopath := os.Getenv("GOPATH")
//
// 	exportfile := ""
// 	if funcType == FUNC_TYPE_HANDLE {
// 		exportfile = fmt.Sprintf("%s/src/%s/handle/mock.go", gopath, module)
// 		os.MkdirAll(path.Dir(exportfile), os.ModeDir)
// 	} else if funcType == FUNC_TYPE_PROXY {
// 		exportfile = fmt.Sprintf("%s/src/%s/proxy/mock.go", gopath, module)
// 		os.MkdirAll(path.Dir(exportfile), os.ModeDir)
// 	} else {
// 		return
// 	}
//
// 	exist, _ := PathExists(exportfile)
// 	if !exist {
// 		AppendFile(exportfile, "package handle\n\n")
// 	}
// 	AppendFile(exportfile, code)
// }
