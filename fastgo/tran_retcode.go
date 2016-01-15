package main
import (
    "strings"
    "fmt"
    "errors"
)


func procData(srcData string) (string, error) {
    newErr := errors.New("constStr invalid")

    sepCheck := ""
    if strings.Contains(srcData, "const (") {
        sepCheck = "const ("
    } else {
        sepCheck = "const("
    }

    strList1 := strings.Split(srcData, sepCheck)
    if len(strList1) < 2 {
        return "", newErr
    }

    strList2 := strings.Split(strList1[1], ")")
    if len(strList2) < 2 {
        return "", newErr
    }

    retStr := srcData

    newConstStr := fmt.Sprintf("var constStr = `const (%s)`\n", strList2[0])
    retStr += newConstStr

    retStr += RETCODE_FUNC_INIT_STR
    return retStr, nil
}

const RETCODE_FUNC_INIT_STR  = `
func init() {
    ret, err := constToMap(constStr)
    if err != nil {
        panic("constStr invalid")
    }
    retcodeMap = ret
}
`