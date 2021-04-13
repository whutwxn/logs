package logs

import (
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

/*
@Description: log分割
@Author : weixiaonan
@Time : 2019/8/12 21:05
*/
var (
	logs = make(map[string]*log.Logger)
	maxAge time.Duration
	rotationTime time.Duration
	logPath string
)


func InitLogs(path string) {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)
	logPath = path
	maxAge = 7*24*time.Hour
	rotationTime = 24*time.Hour
}

func addNewLogFile(logName string) {
	l:=log.New()
	l.SetOutput(os.Stdout)
	l.AddHook(newRotateHook(logPath, logName, maxAge, rotationTime))
	logs[logName]=l
}

func newRotateHook(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) *lfshook.LfsHook {
	baseLogPath := path.Join(logPath, "%Y-%m-%d--"+logFileName+".log")
	writer, err := rotatelogs.New(
		baseLogPath,
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	return lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
}

func checkLog(logName string,args ...interface{}) bool {
	if logs[logName] == nil {
		addNewLogFile(logName)
	}
	if len(args) <= 0 {
		return false
	}
	return true
}

func Info(logName string,args ...interface{}) {
	if !checkLog(logName,args){
		return
	}
	logs[logName].Info(args, printCallerName())
}

func Debug(logName string,args ...interface{}) {
	if !checkLog(logName,args){
		return
	}
	logs[logName].Debug(args, printCallerName())
}

func Warn(logName string,args ...interface{}) {
	if !checkLog(logName,args){
		return
	}
	logs[logName].Warn(args, printCallerName())
}

func Error(logName string,args ...interface{}) {
	if !checkLog(logName,args){
		return
	}
	logs[logName].Error(args, printCallerName())
}

func Fatal(logName string,args ...interface{}) {
	if !checkLog(logName,args){
		return
	}
	logs[logName].Fatal(args, printCallerName())
}

func Panic(logName string,args ...interface{}) {
	if !checkLog(logName,args){
		return
	}
	logs[logName].Panic(args, printCallerName())
}

func printCallerName() string {
	pc, file, line, _ := runtime.Caller(2)
	lineNum := strconv.Itoa(line)
	files := strings.Split(file, "/")
	pcs := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	return " -->path: [func:" + pcs[len(pcs)-1] + "(" + files[len(files)-1] + ":" + lineNum + ")]"
}

func LogLevel(level string) error {
	if level == "debug" {
		for _, logger := range logs {
			logger.SetLevel(log.DebugLevel)
		}
	} else if level == "info" {
		for _, logger := range logs {
			logger.SetLevel(log.InfoLevel)
		}
	}else if level == "error" {
		for _, logger := range logs {
			logger.SetLevel(log.ErrorLevel)
		}
	}else {
		return errors.New("level error")
	}
	return nil
}

