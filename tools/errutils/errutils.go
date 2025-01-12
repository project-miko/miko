package errutils

import (
	"fmt"
	"runtime"
)

func PrintPanicStackError() {
	if x := recover(); x != nil {
		fmt.Println("panic ", x)
		PrintPanicStack()
	}
}

func PrintPanicStack() {
	for i := 0; i < 10; i++ {
		funcName, file, line, ok := runtime.Caller(i)
		if ok {
			funcName := runtime.FuncForPC(funcName).Name()
			errInfo := fmt.Sprintf("frame %d:[func:%s, file: %s, line:%d]", i, funcName, file, line)
			fmt.Println(errInfo)
		}
	}
}
