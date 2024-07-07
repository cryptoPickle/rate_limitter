package main

import (
	"fmt"
	"time"
)

func main() {
	duration := time.Second * 2

	ticker := time.NewTicker(duration)

	for {
		<-ticker.C
		fmt.Println("tick")
	}
}
