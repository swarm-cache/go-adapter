package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/swarm-cache/go-adapter/export"
)

func main() {

	err, c := export.Connect("127.0.0.1:3666")
	if err != nil {
		fmt.Println(err)
	}
	i := 0

	for i < 1000 {
		str := "love" + strconv.Itoa(i)
		bStr := []byte(str)

		if err := c.Set(str, &bStr); err != nil {
			fmt.Printf("Error! %s\n", err)
		}

		i++
	}

	fmt.Println(c.Get("love10"))
	fmt.Println(c.Del("love10"))
	fmt.Println(c.Get("love10"))

	time.Sleep(5 * time.Second)
}
