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

func Log(msg error) {
	if !SendLogs {
		return
	}

	currentTime := time.Now().Local()
	logMsg := fmt.Sprintf("%s: %v\n", currentTime.String(), msg)

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

func CreateDetails(logInfo string, errMsg error) error {
	_, _, line, _ := runtime.Caller(1)
	return fmt.Errorf("%s %w. line: %d", logInfo, errMsg, line)
}
