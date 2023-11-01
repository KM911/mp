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
	fmt.Println("error exit with code ", err.Type)
	switch err.Type {
	default:
		fmt.Println(err.Msg)
	}
}

var (
	APP_NAME = "mp" // manage profile
	Export   = false
	HELP     = "mp is a tool for managing environment path."
	USAGE    = `Usage: 
    pwd | mp             add current folder into env PATH
    mp  [value]          add value into env PATH
    mp  [key] [value]    add env key=value`
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
	fmt.Println(HELP)
	fmt.Println()
	fmt.Println(USAGE)
	fmt.Println()
}

func Reminder() {
	time.Sleep(1 * time.Second)
	HelpInfo()
	CleanPath()
	os.Exit(0)
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
}
