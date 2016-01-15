package main
import (
    "strings"
    "strconv"
    "errors"
)


const(
    RET_OK = 0          // ok
    RET_ERROR = 9999    // fail
)


func GetRetMsg(retcode int) string {
    retMsg, ok := retcodeMap[retcode]
    if !ok {
        return ""
    }

    return retMsg
}


//=================================================
var retcodeMap map[int]string

func constToMap(constStrInput string) (map[int]string, error) {
    newErr := errors.New("constStr invalid")

    strList1 := strings.Split(constStrInput, "const")
    if len(strList1) < 2 {
        return nil, newErr
    }

    strList2 := strings.Split(strList1[1], "(")
    if len(strList2) < 2 {
        return nil, newErr
    }

    strList3 := strings.Split(strList2[1], ")")
    if len(strList3) < 2 {
        return nil, newErr
    }

    strList4 := strings.Split(strings.TrimSpace(strList3[0]), "\n")
    if len(strList4) < 2 {
        return nil, newErr
    }

    retMap := map[int]string{}
    for _, row := range strList4 {
        if len(strings.TrimSpace(row)) == 0 {
            continue
        }

        rowStrList1 := strings.Split(strings.TrimSpace(row), "=")
        if len(rowStrList1) < 2 {
            return nil, newErr
        }

        rowStrList2 := strings.Split(rowStrList1[1], "//")
        if len(rowStrList2) < 2 {
            return nil, newErr
        }

        retcode, _ := strconv.Atoi(strings.TrimSpace(rowStrList2[0]))
        retmsg := strings.TrimSpace(rowStrList2[1])
        retMap[retcode] = retmsg
    }

    return retMap, nil
}