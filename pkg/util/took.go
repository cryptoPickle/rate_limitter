package util

import (
	"fmt"
	"time"
)

type returnFunc func()

func Took(name string) returnFunc {
	start := time.Now()
	return func() {
		took := time.Since(start)
		fmt.Printf("%v took %v\n", name, took)
	}
}
