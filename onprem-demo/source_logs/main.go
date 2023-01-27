package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/nxadm/tail"
)

func main() {
	httpEndpoint := os.Getenv("HTTP_ENDPOINT")
	logFile := os.Getenv("LOG_FILE")
	// check if log file exists and panic if not
	_, err := os.Stat(logFile)
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("Log file does not exist: %s", logFile))
	}

	fmt.Printf("httpEndpoint --------------------------------------\n")
	spew.Dump(httpEndpoint)

	t, err := tail.TailFile(
		logFile, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		panic(err)
	}

	counter := 0

	// Print the text of each received line
	for line := range t.Lines {
		fmt.Println(line.Text)
		counter++

		if counter == 10 {
			break
		}
	}
}
