package main

import (
	"flag"
	"fmt"
	"github.com/KM911/oslib/adt"
	"github.com/KM911/oslib/fs"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

type ErrorMessage struct {
	Type int
	Msg  string
}

func FileLogger(src string) {

	logFile, err := os.OpenFile(src, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
}

func EmitError(n int, msg string) {
	panic(ErrorMessage{
		Type: n,
		Msg:  msg,
	})
}

func Recover(errorHandler func(ErrorMessage)) {
	if errMsg := recover(); errMsg != nil {
		errorHandler(errMsg.(ErrorMessage))
		log.Println(string(debug.Stack()))
	}
}

// TODO add error message
func ErrorHandler(err ErrorMessage) {
	switch err.Type {
	default:
		fmt.Println(err.Type, err.Msg)
	}
}

var (
	APP_NAME = "mp" // manage profile
	Export   = false
)

func IsInPath(src string) bool {
	index := strings.Index(os.Getenv("PATH"), src)
	if index == -1 {
		return false
	}
	return true
}

func ParseUserInput() {
	ArgsLens := len(flag.Args())
	//fmt.Println("args lens is ", ArgsLens)
	//fmt.Println("args is ", flag.Args())
	// 关于 flag的解析 问题
	switch ArgsLens {
	case 0:
		ModifyEnvironment(CheckPath(ReadFromStream()))
	case 1:
		ModifyEnvironment(CheckPath(flag.Arg(0)))
	case 2:
		ModifyEnvironment(CheckValue(flag.Arg(0), flag.Arg(1)))
	default:
		HelpInfo()
	}
}

func HelpInfo() {
	fmt.Println(APP_NAME + " useage is \n")
	fmt.Println("Add current folder to environment path")
	fmt.Println("	pwd | " + APP_NAME + "    add path by pipeline")
	fmt.Println("Add path by args")
	fmt.Println("	", APP_NAME+" [path]")
	fmt.Println("	", APP_NAME+" [key] [value]")
}
func Reminder() {
	time.Sleep(1 * time.Second)
	HelpInfo()
	os.Exit(0)
	//EmitError(3, "without pipe")
}

func init() {
	// TODO add init
	//flag.BoolVar(&Export, "export", false, "export to system")
	flag.BoolVar(&Export, "e", false, "export to system")
	flag.Parse()
}

func main() {
	FileLogger(filepath.Join(fs.ExecutePath(), "error.log"))
	defer adt.TimerStart().End()
	defer Recover(ErrorHandler)
	ParseUserInput()
	//fmt.Println(QueryUserPath())
}
