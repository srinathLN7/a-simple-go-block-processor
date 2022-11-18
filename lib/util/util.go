package util

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	LOG         string = "[main] 				"
	NODE_LOG    string = "[nodeService-routine] "
	PROCESS_LOG string = "[processor-routine]   "
	BC_LOG      string = "[bc]   				"
)

// GetCurrentTimeStamp:: returns the current timestamp in GMT+0 zone
// A standard time zone is conisdered to avoid any potential effects of non-determinism
func GetCurrentTimeStamp() string {
	currentTime := time.Now()
	currentTimeStr := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second())
	return (currentTimeStr + ".000Z")
}

//GetRandomNum:: generates a random number with in the certain range
func GetRandomNum() int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 10

	return rand.Intn(max-min+1) + min
}

//GetLogStr:: prints `*` up to the width of the terminal
func GetLogStr() string {
	var logStr string
	width, _, _ := terminal.GetSize(0)
	for i := 0; i < width-20; i++ {
		logStr = logStr + "*"
	}
	return logStr
}
