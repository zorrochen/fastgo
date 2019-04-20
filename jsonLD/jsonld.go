package jsonLD

import (
	"encoding/json"
	"reflect"
)

/**
  类json-LD格式，生成字段说明
*/
func Unmarshal(s string) (data string, linked map[string]string) {
	decodeObj := map[string]interface{}{}
	err := json.Unmarshal([]byte(s), &decodeObj)
	if err != nil {
		return
	}

	linkedStr := map[string]string{}
	linkedObj, ok := decodeObj["@context"]
	if !ok {
		return s, linkedStr
	}

	if reflect.TypeOf(linkedObj).Kind() != reflect.Map {
		return
	}
	for k, v := range linkedObj.(map[string]interface{}) {
		if reflect.TypeOf(v).Kind() != reflect.String {
			continue
		}
		linkedStr[k] = v.(string)
	}

	delete(decodeObj, "@context")
	dataJson, _ := json.Marshal(decodeObj)
	return string(dataJson), linkedStr
}
