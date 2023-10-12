package main

import (
	"fmt"
	oslib "github.com/KM911/oslib/cmd"
	"github.com/KM911/oslib/fc"
	"github.com/KM911/oslib/fs"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// TODO add Windows version
var (
	USERregpath    = "HKEY_CURRENT_USER\\Environment"
	MACHINEregpath = "HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet" +
		"\\Control\\Session Manager\\Environment"
	ENVIRONMENT = map[string]string{
		"PATH": "",
	}
	StringBuilder strings.Builder
	USERPATH      = QueryUserPath()
)

// 修改reg来修改环境变量
func ModifyEnvironment(data string) {
	SetUserPath(data)
	//fmt.Println("Modify argv is", data)
}

//func IsInPath(src string) bool {
//	// TODO 判读是否存在的环境变量
//	index := strings.Index(QueryUserPath(), src)
//	if index == -1 {
//		return false
//	}
//	return true
//}

// 我们默认也是添加到用户的环境变量中 这样主要是可以避免权限问题

func CheckPath(src string) string {
	src = filepath.ToSlash(src)
	src = strings.ToUpper(src)
	if Export {
		ExportPath(src)
		os.Exit(1)
	}
	if !fs.IsExit(src) {
		fmt.Println(src, " is not exits")
		os.Exit(2)
	}
	if strings.HasSuffix(src, ".mp") || strings.HasSuffix(src, ".MP") {
		LoadFromFile(src)
		os.Exit(0)
	}
	if IsInPath(src) {
		// do not need to add
		EmitError(4, src+" is in path")
		os.Exit(4)
	}
	return src
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
	src = strings.ToUpper(strings.TrimSpace(string(all)))
	fmt.Println("pipe is ", src)
	return src

}
func SetUserPath(value string) {
	userPath := QueryUserVariable("PATH")
	if IsInPath(value) {
		return
	}
	pc := 0
	for {
		if pc == len(userPath) || userPath[pc] != 59 {
			break
		} else {
			pc++
		}
	}
	userPath = userPath[pc:]
	fmt.Println("user path is", userPath)
	if !strings.HasSuffix(userPath, ";") {
		userPath = userPath + ";"
	}
	userPath = userPath + value
	SetUserVaiable("PATH", userPath)
}

func SetUserVaiable(k, v string) {
	command := fmt.Sprint("reg add HKEY_CURRENT_USER\\Environment"+
		" /v ", k, " /t REG_SZ /d ", v, " /f")
	oslib.RunStd(command)
}

func SetSystemPath(value string) {

}

func SetSystemVariable(k, v string) {
}

func QueryUserPath() string {
	userPath := QueryUserVariable("PATH")
	if strings.HasSuffix(userPath, ";") {
		return userPath
	} else {
		return userPath + ";"
	}
}

func QueryUserVariable(k string) string {
	command := "reg query " + USERregpath + " /v " + k
	value := oslib.RunReturn(command)
	word_list := strings.Split(value, "   ")
	if len(word_list) < 3 {
		fmt.Println(k, "is not set")
		return ""
	}
	word := strings.TrimSpace(strings.Split(value, "   ")[3])
	return strings.ToUpper(word)
}

func QuerySystemPath() string {
	return QuerySystemVariable("PATH")
}

func QuerySystemVariable(k string) string {
	value := oslib.RunReturn(
		"reg query " + MACHINEregpath + " /v" +
			" " + k + "")
	word_list := strings.Split(value, "   ")
	if len(word_list) < 3 {
		fmt.Println("mp is not set")
		return ""
	}
	word := strings.TrimSpace(strings.Split(value, "   ")[3])
	return word
}

//
func LoadFromFile(src string) {
	LoadFile(src)
	//fmt.Println("env is", ENVIRONMENT)
	for s, s2 := range ENVIRONMENT {
		if s == "PATH" {
			// 一次性将全部的path添加了 导致判断失效了
			SetUserPath(s2)
		} else {
			SetUserVaiable(s, s2)
		}
	}
	fmt.Println("load from file success")
}

func LoadFile(src string) {
	// 直接将文件中的内容全部读取就好了 想复杂了不是吗 笑死我了
	file, err := os.Open(src)
	if err != nil {
		return
	}
	stat, err := file.Stat()
	if err != nil {
		return
	}
	buffer := make([]byte, stat.Size())
	_, err = file.Read(buffer)
	if err != nil {
		return
	}
	defer file.Close()
	lines := strings.Split(strings.TrimSpace(string(buffer)), "\n")
	filted_lines := fc.Filter(lines, IsEmptyLine)
	for i := range filted_lines {
		ParseFileContent(filted_lines[i])
	}
	// 原始 + 现在的

	ENVIRONMENT["PATH"] = USERPATH + StringBuilder.String()
	return

}
func IsEmptyLine(line string) bool {
	return len(strings.TrimSpace(line)) != 0
}
func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// export user path for backup
func ExportPath(src string) {
	create, err := os.Create(src)
	if err != nil {
		return
	}
	defer create.Close()

	UserPath := strings.Split(QueryUserPath(), ";")
	sort.Strings(UserPath)
	// 删除重复的路径
	UserPath = removeDuplicates(UserPath)
	for _, s := range UserPath {
		create.WriteString(s)
		create.WriteString("\n")
	}
	// 导出正常的用户环境变量
	value := oslib.RunReturn(
		"reg query HKEY_CURRENT_USER\\Environment")
	words := strings.Split(value, "   ")
	lens := len(words)
	for i := 1; i < lens; i += 3 {
		if strings.TrimSpace(words[i]) == "Path" {
			continue
		} else {
			create.WriteString(strings.TrimSpace(words[i]))
			create.WriteString("=")
			create.WriteString(strings.TrimSpace(words[i+2]))
			create.WriteString("\n")
		}
	}
}

func ParseFileContent(line string) {
	line = strings.TrimSpace(strings.ToUpper(line))
	equalIndex := strings.Index(line, "=")
	if equalIndex == -1 {
		if IsInPath(line) {
			return
		}
		StringBuilder.WriteString(line)
		StringBuilder.WriteString(";")
	} else {
		k := line[:equalIndex]
		v := line[equalIndex+1:]
		ENVIRONMENT[k] = strings.TrimSpace(v)
	}
}
