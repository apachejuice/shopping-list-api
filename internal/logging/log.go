package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/palantir/stacktrace"
)

var (
	logFile string
	handle  *os.File
	logger  *log.Logger
)

func InitLog(file string) {
	logFile = file

	// Destroy if the previous file exists, create it if not
	stat, err := os.Lstat(logFile)
	if err != nil {
		log.Fatal(err)
	} else if os.IsExist(err) {
		log.Println("Destroying previous log file")
	} else if os.IsNotExist(err) {
		os.Create(logFile)
	}
	if stat.IsDir() {
		log.Fatalf("Log file %q is a directory\n", logFile)
	}

	handle, err = os.OpenFile(logFile, os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}

	logger = log.New(handle, "[shapi] ", log.LstdFlags)
	logger.Printf("Hello! Starting up at %s\n", time.Now().UTC())
	gin.DefaultWriter = logger.Writer()
	gin.DefaultErrorWriter = logger.Writer()
}

func Info(msg string, args ...any) {
	pc, filename, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	fileParts := strings.Split(filename, "/")
	file := fileParts[len(fileParts)-1]

	nameParts := strings.Split(fn.Name(), "/")
	fnname := nameParts[len(nameParts)-1]

	logger.Printf(fmt.Sprintf("%s() at %s:%d: ", fnname, file, line)+msg+"\n", args...)
}

func Error(err error, code string) {
	logger.Println("ERROR: " + code)
	logger.Println("Root cause: " + stacktrace.RootCause(err).Error() + "\n")
	logger.Println(err.Error())
}
