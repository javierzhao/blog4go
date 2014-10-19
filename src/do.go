package main

import (
	"fmt"
	"strings"
	"utils"
)

var logChannel chan string = make(chan string)

func loggingLoop() {
	for {
		//wait for a message to arrive
		msg := <-logChannel

		//log the msg
		fmt.Println(msg)
	}
}

func main() {
	// go loggingLoop()

	//do some stuff here
	// logChannel <- "messaged to be logged"
	//do other stuff here
	//
	tagNames := "a,bc"

	li := strings.LastIndex(tagNames, ",")
	l := len(tagNames)
	if li != -1 && l > 0 {
		tagNames = utils.Substr(tagNames, 0, li)
	}

	fmt.Println(tagNames)
}
