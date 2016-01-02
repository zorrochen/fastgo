package main

import (
	"fmt"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"runtime"
	"log"
)

var (
	showVersion  = flag.Bool("v", false, "print version string")

	srcFileDir  = flag.String("srcFileDir", "./srcFiles", "srcFileDir")
)



func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println("v1.0")
		return
	}

	// 设置使用的CPU个数
	log.Print("start...(CPU:%d)", runtime.NumCPU())

	// 信号处理
	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	signalIgnoreChan := make(chan os.Signal, 1)
	go func() {
		for {
			select {
			case <-signalChan:
				exitChan <- 1
				return
			case <-signalIgnoreChan:
				continue
			}
		}
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(signalIgnoreChan, syscall.SIGPIPE)
	
	server := NewServer()
	server.Start()

	<-exitChan
	log.Print("server stoped")
	server.Exit()
}
