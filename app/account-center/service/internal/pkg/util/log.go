package util

import (
	"fmt"
	"runtime"
	"time"
)

func Println(msg ...interface{}) {
	fmt.Println(append([]interface{}{time.Now().Format("2006-01-02 15:04:05")}, msg...)...)
}

func Stack(err interface{}) {
	fmt.Printf("%v\n%s", err, stack())
}

func stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], true)])
}

