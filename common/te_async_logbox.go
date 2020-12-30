package common

import (
	"os"
	"strings"
	"time"
)

type RunningLog struct {
	LogLevel []string
	LogFile  string
	LogQueue []struct{
		Level   int
		Message string
	}
}


func CreateLog() *RunningLog {
	res := &RunningLog{
		LogFile: "log/access.log",
		LogLevel: []string{"ACCESS", "INFO", "WARN", "ERROR"},
	}

	go res.syncToDisk()
	return res
}


func (rl *RunningLog) Add(level int, message string) {
	rl.LogQueue = append(rl.LogQueue, struct {
		Level   int
		Message string
	}{Level: level, Message: message})
}


func (rl *RunningLog) syncToDisk() {
	fd, _ := os.OpenFile(rl.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	for Running {
		if rl.LogQueue != nil {
			for _, logStruct := range rl.LogQueue {
				fdTime := time.Now().Format("2006-01-02 15:04:05")
				fdContent := strings.Join([]string{fdTime, " ==> ", rl.LogLevel[logStruct.Level]," ==> ", logStruct.Message, "\n"}, "")
				fd.Write([]byte(fdContent))
				rl.LogQueue = rl.LogQueue[1:]
			}
		}
		time.Sleep(1 * time.Second)
	}
	fd.Close()
}