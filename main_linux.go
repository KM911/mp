package main

import (
	"fmt"
	"github.com/KM911/oslib/fs"
	"io"
	"os"
	"strings"
)

// check the path is valid and return command
// linux path is case-sensitive
var (
	SRC      = "~/.profile"
	PathList = map[string][]string{}
)

func CheckPath(src string) string {
	if !fs.IsExit(src) {
		EmitError(2, "path is not exits")
		os.Exit(-1)
	}
	if IsInPath(src) {
		EmitError(4, "path is in path , do not need to add")
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
	file, err := os.OpenFile(SRC, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = fmt.Fprintf(file, data)
	if err != nil {
		EmitError(1, "write file error")
	}
	fmt.Println("add value success")
}

func QueryUserPath() string {
	return os.Getenv("PATH")
}

func LoadProfile() {
	data, err := os.ReadFile("~/.profile")
	if err != nil {
		EmitError(1, "read file error")
	}
	//
	lines := strings.Split(string(data), "\n")
	for i := range lines {
		if strings.HasPrefix(lines[i], "export") {
			///mnt/d/CODE/go/mp
			exports := strings.Split(lines[i][8:], "=")
			src := ""
			if exports[0] == "PATH" {
				src = strings.Split(exports[1], ":")[0]
			} else {
				src = exports[1]
			}
			if fs.IsExit(src) {
				PathList[exports[0]] = append(PathList[exports[0]], src)
			}

		}
	}
}

func CleanPath() {
	fmt.Println("Not implement clean path")
	// paths := strings.Split(QueryUserPath(), ":")
	// //fmt.Println(paths)
	// //paths = paths[:len(paths)-1]
	// //// 在一个循环中删除元素应该使用迭代器
	// for i := 0; i < len(paths); i++ {
	// 	if !fs.IsExit(paths[i]) {
	// 		fmt.Println("remove not exit path :", paths[i])
	// 		paths = append(paths[:i], paths[i+1:]...)
	// 		i--
	// 	}
	// }
	// /// 删除重复元素
	// paths = removeDuplicates(paths)
	// //
	// //// 将path的值排序
	// sort.Strings(paths)
	// //fmt.Println("Current path is ")
	// //for i := range paths {
	// //	fmt.Println(paths[i])
	// //}
	// //SetUserVaiable("PATH", strings.ReplaceAll(strings.Join(paths, ";"), "/", `\`)+";")
}
