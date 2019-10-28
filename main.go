package main

import (
	"fmt"
	"quickstart/lib"
	"time"
)

func main() {
	timestamp := time.Now().UnixNano()
	fmt.Println(timestamp)
	for i := 0; i < 100000; i++ {
		lib.GetIDInstance().Get()
	}
	timeDiff := (time.Now().UnixNano() - timestamp) * 1000 / 1000000000

	s := fmt.Sprintf("time use: %d", timeDiff)
	fmt.Println(s)
}
