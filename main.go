package main

import (
	. "fastgo/handler"
	"flag"
	"log"
	"runtime"
)

var (
	serviceName = flag.String("srv", "tmp", "service")
	filepath    = flag.String("filepath", "", "filepath")
	funcType    = flag.Int("type", 1, "function type")
	mockFlag    = flag.Bool("mock", false, "mock data switch")
)

func main() {
	flag.Parse()

	// 设置使用的CPU个数
	log.Printf("start...(CPU:%d)", runtime.NumCPU())

	info, err := ReadFile(*filepath)
	if err != nil {
		return
	}
	genRstFormat := ParseAndGen(info)
	Export(genRstFormat, *serviceName, "newexport")

	log.Print("server stoped")
}
