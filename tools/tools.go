package tools

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

func GetRandChar() uint8 {
	chars := "0123456789abcdefghijklmnopqrstuvwxyz"
	clen := len(chars)

	pos := rand.Int() % clen

	return chars[pos]
}

func GetRandStr(amount int) string {
	buf := new(strings.Builder)

	var i int = 0
	var c uint8
	for ; i < amount; i++ {
		c = GetRandChar()
		buf.WriteByte(c)
	}

	return buf.String()
}

// get the caller of the function
func GetCaller(skip int) string {
	_, file, line, _ := runtime.Caller(skip + 1)
	file = file[strings.LastIndex(file, "/")+1:]
	return fmt.Sprintf("%s:%d", file, line)
}

// get milliseconds
func GetMillisecond(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func FormatShortData(target time.Time) (time.Time, error) {
	layout := "2006.01.02"
	timeStr := target.Format(layout)
	result, err := time.Parse(layout, timeStr)
	return result, err
}

func IsPathExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {

		return os.IsExist(err)
	}

	return true
}

func Int16ToBytes(x int16) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
