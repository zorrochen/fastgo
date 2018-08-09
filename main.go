package main

import (
	. "fastgo/handler"
	"flag"
	"log"
	"path"
	"runtime"
)

var (
	serviceName = flag.String("srv", "tmp", "service")
	filepath    = flag.String("filepath", "", "filepath")
	testFlag    = flag.Bool("t", false, "test switch")
)

func main() {
	flag.Parse()

	// 启动打印(附:cpu个数)
	log.Printf("start...(CPU:%d)", runtime.NumCPU())

	info, err := ReadFile(*filepath)
	if err != nil {
		return
	}
	genRstFormat := ParseAndGen(info)
	funcName := path.Base(*filepath)
	Export(genRstFormat, *serviceName, funcName)

	if *testFlag {
		ExportTests(*serviceName, funcName)
	}

	log.Print("server stoped")
}
