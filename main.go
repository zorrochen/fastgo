package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/zorrochen/fastgo/handler"
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

	info, err := handler.ReadFile(*filepath)
	if err != nil {
		return
	}

	// code
	genRstFormat := handler.ParseAndGen(info)
	handler.ExportInCurrentPath(genRstFormat, *filepath)

	// testcase
	if *testFlag {
		handler.ExportTestsInCurrentPath(*filepath)
	}

	log.Print("server stoped")
}
