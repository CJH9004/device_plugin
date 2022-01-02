package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	for {
		fmt.Println(">>>>>>>>>> availiable devices:")
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, "DEVICE_ID") {
				fmt.Println(env)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
