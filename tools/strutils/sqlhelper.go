package strutils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/project-miko/miko/tools/log"
)

func Prefix0x(data string) string {
	if strings.HasPrefix(data, "0x") || strings.HasPrefix(data, "0X") {
		return data
	}
	return "0x" + data
}

func Trim0x(data string) string {
	if strings.HasPrefix(data, "0x") || strings.HasPrefix(data, "0X") {
		data = data[2:]
	}
	return data
}

func Trim0xAndToLower(data string) string {
	data = strings.ToLower(data)
	data = Trim0x(data)
	return data
}

func StringSliceToInString(s []string) string {
	if len(s) == 0 {
		return "''"
	}

	var builder strings.Builder

	for _, v := range s {
		if builder.Len() > 0 {
			builder.WriteRune(',')
		}
		builder.WriteRune('\'')
		builder.WriteString(v)
		builder.WriteRune('\'')
	}

	return builder.String()
}

// receives array like []int64{1,2,3...}, returns string like '1','2','3'...
func IdsToInString(idArr []int64) string {
	if len(idArr) == 0 {
		return "''"
	}

	var idBytes []byte

	for _, v := range idArr {
		if len(idBytes) > 0 {
			idBytes = strconv.AppendQuoteRune(idBytes, ',')
		}

		idBytes = strconv.AppendInt(idBytes, v, 10)
	}

	return "'" + string(idBytes) + "'"
}

func FormatNoticeMsg(content, params string) string {

	paramMap := make(map[string]interface{})
	e := json.Unmarshal([]byte(params), &paramMap)
	if e != nil {
		log.Error("", "FormatNoticeMsg failed %s %s %s", content, params, e.Error())
		return ""
	}

	for k, v := range paramMap {
		str := fmt.Sprintf("%v", v)

		content = strings.ReplaceAll(content, k, str)
	}

	return content
}
