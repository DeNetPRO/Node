package logger

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

var (
	SendLogs = true
)

// ====================================================================================

func LogError(logInfo string, errMsg error) {
	if !SendLogs {
		return
	}

	currentTime := time.Now().Local()
	logMsg := fmt.Sprintf("%s: %s: %v\n", currentTime.String(), logInfo, errMsg)

	fmt.Println(logMsg)

	url := "http://68.183.215.241:9091/logs/node/"

	req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte(logMsg)))

	client := &http.Client{Timeout: time.Minute}

	_, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
}

// ====================================================================================

func GetDetailedError(errMsg error) error {
	_, _, line, _ := runtime.Caller(1)
	return fmt.Errorf("%w. line: %d", errMsg, line)
}
