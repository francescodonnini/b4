package logging

import (
	"b4/shared"
	"fmt"
	"log"
	"os"
	"sync"
)

var lock = &sync.RWMutex{}

var loggers map[string]*log.Logger

var nodeName string

func GetInstance(name string) *log.Logger {
	if loggers == nil {
		lock.Lock()
		defer lock.Unlock()
		if loggers == nil {
			loggers = make(map[string]*log.Logger)
			ip := shared.GetAddress()
			nodeName = ip.Ip
		}
	}
	logger, ok := loggers[name]
	if !ok {
		logFile, err := os.OpenFile(fmt.Sprintf("/var/%s-%s", nodeName, name), os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Panicln(err)
		}
		logger = log.New(logFile, "", log.Lshortfile|log.LstdFlags)
		loggers[name] = logger
	}
	return logger
}
