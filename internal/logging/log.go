package logging

import (
	"log"
	"os"

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

	logger = log.New(handle, "[shapi] ", log.LstdFlags|log.Lshortfile)
	gin.DefaultWriter = logger.Writer()
	gin.DefaultErrorWriter = logger.Writer()
}

func Info(msg string, args ...any) {
	logger.Printf(msg+"\n", args...)
}

func Error(err error, code string) {
	logger.Println("ERROR: " + code)
	logger.Println("Root cause: " + stacktrace.RootCause(err).Error() + "\n")
	logger.Println(err.Error())
}
