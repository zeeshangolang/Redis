package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RedisFuncs func(string) func()

// var
// Update Commands will
var mu sync.Mutex

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	Redis := make(map[string]any)

	for scanner.Scan() {
		fmt.Print("\n  ")
		TEXT := scanner.Text()
		Mtext := strings.Fields(TEXT)

		AllFunc(Mtext, Redis)
	}
	if err := scanner.Err(); err != nil {
		fmt.Print(err)
	}

}

func AllFunc(InputValues []string, MAp map[string]any) any {

	var ReturnVal any

	length := len(InputValues)

	switch {

	case length == 2 && InputValues[0] == "Get":
		ReturnVal = GetFunc(InputValues[1], MAp)
		fmt.Print(ReturnVal, "\n")
		return ReturnVal

	case length == 3 && InputValues[0] == "Set":
		MAp[InputValues[1]] = InputValues[2]
		ReturnVal = "Ok\n"
		fmt.Print(ReturnVal, "\n")

	case length <= 2 && InputValues[0] == "Ping":
		if length == 1 {
			fmt.Print("Pong\n")

		} else {
			fmt.Print("Pong ", InputValues[1], "\n")
		}

	case length > 1 && InputValues[0] == "Del":

		DeleteKeys(MAp, InputValues[1:])

	case length == 3 && InputValues[0] == "Expire":
		ExpireInt, err := strconv.Atoi(InputValues[2])
		if err != nil {
			fmt.Print("please use correct values for expiry")
		}

		val := AddEXpiry(InputValues[1], ExpireInt, MAp)
		fmt.Print(val)
	default:
		fmt.Print("no command found\n")
	}

	return ReturnVal

}

func AddEXpiry(val string, Expiry int, RedisMap map[string]any) int {
	timeLimit := time.Now().Add(time.Duration(Expiry) * time.Second)
	TimeDuration := time.Until(timeLimit)
	time.AfterFunc(TimeDuration, func() {
		mu.Lock()
		delete(RedisMap, val)
		mu.Unlock()
	})
	if _, ok := RedisMap[val]; ok {
		return 1
	}
	return 0
}

func DeleteKeys(Mape map[string]any, DelVals []string) {

	for _, val := range DelVals {

		delete(Mape, val)
	}
}

func GetFunc(key string, PMap map[string]any) any {
	return PMap[key]
}
