package handler

import (
	"fmt"
	"github.com/cweill/gotests"
	"strings"
)

func gotestsRun(fileFullName string) {
	exclRE, err := parseRegexp("^make*")
	if err != nil {
		return
	}

	gts, err := gotests.GenerateTests(fileFullName, &gotests.Options{Exclude: exclRE})
	if err != nil {
		return
	}

	if len(gts) == 0 {
		return
	}

	allTestsCode := ""
	for _, v := range gts {
		allTestsCode += string(v.Output)
	}

	testFileFullName := strings.Replace(fileFullName, ".go", "_test.go", -1)
	exist, _ := PathExists(testFileFullName)
	// 已存在，不写入
	if exist {
		fmt.Printf("already exist.\n")
		return
	}
	writeFile(testFileFullName, allTestsCode)
}
