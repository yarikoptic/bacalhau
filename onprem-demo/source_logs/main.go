package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/nxadm/tail"
)

type StreamingResult struct {
	// one of these two will be defined
	// LocalPath for images on the filesystem
	LocalPath string `json:"LocalPath"`
	// InlineData for things like log lines
	InlineData string `json:"InlineData"`
	// the Channel to broadcast the streaming result on
	// this is a logical channel so we can use the same
	// gossip sub topic otherwise we have to manually connect
	Channel string `json:"Channel"`
}

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

	// Print the text of each received line
	for line := range t.Lines {
		fmt.Println(line.Text)

		if strings.Contains(line.Text, "DHCPACK") {
			r := StreamingResult{
				InlineData: line.Text,
				Channel:    "ap_connections",
			}
			bs, err := json.Marshal(r)
			if err != nil {
				log.Printf("err yielding result: %s", err)
				continue
			}
			buf := bytes.NewReader(bs)
			resp, err := http.Post("http://172.17.0.1:9600/publish", "application/json", buf)
			if err != nil {
				log.Printf("err yielding result: %s", err)
				continue
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("err reading body: %s", err)
				continue
			}
			log.Printf("resp from bacalhau streaming: %s", string(body))
		}
	}
}
