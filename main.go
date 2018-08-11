package main

import (
	"flag"
	. "github.com/zorrochen/fastgo/handler"
	"log"
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

	// code
	genRstFormat := ParseAndGen(info)
	ExportInCurrentPath(genRstFormat, *filepath)

	// testcase
	if *testFlag {
		ExportTestsInCurrentPath(*filepath)
	}

	log.Print("server stoped")
}
