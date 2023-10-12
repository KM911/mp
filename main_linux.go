package main

import (
	"fmt"
	"github.com/KM911/oslib/fs"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// check the path is valid and return command
var (
	SRC = ".profile"
)

if !fs.IsExit(src) {
	fmt.Println("path is not exits")
	func CheckPath(src string) string {
	//EmitError(2, "path is not exits")
	os.Exit(-1)
	}
	if IsInPath(src) {
		// do not need to add
		EmitError(4, "path is in path")
		return ""
	}
	return "export PATH=$PATH:" + src + "\n"
}

func CheckValue(k, v string) string {
	return "export " + k + "=" + v + "\n"
}

func ReadFromStream() (src string) {
	go Reminder()
	// io
	all, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err.Error())
		EmitError(1, "read from stdin error")
		return ""
	}
	src = strings.TrimSpace(string(all))
	fmt.Println("pipe is ", src)
	return src

}

// 向profile文件中写入数据
func ModifyEnvironment(data string) {
	SRC = filepath.Join(os.Getenv("HOME"), SRC)
	src, err := os.OpenFile(SRC, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer src.Close()
	_, err = fmt.Fprintf(src, data)
	if err != nil {
		EmitError(1, "write file error")
	}
	fmt.Println("add value success")
}