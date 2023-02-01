package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
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

var latestImageCid string

/*

curl -H "Content-type: application/json" \
--data '' \
-H "Authorization: Bearer xoxb-not-a-real-token-this-will-not-work" \
-X POST https://slack.com/api/chat.postMessage

*/

func postToSlack(text string) error {
	payload := fmt.Sprintf(`{"text": "%s"}`, text)
	http.Post(os.Getenv("SLACK_WEBHOOK_URL"), "application/json", strings.NewReader(payload))
	return nil
}

func notifySlack(labels []string) {
	// make an http request to slack webhook
	postToSlack(fmt.Sprintf("detected %+v!", labels))
}

func processInference(latestImageCid string) {
	// pretend to do AI inference to get labels
	labels := []string{"fish", "dog"}
	notifySlack(labels)
}

func processAPJoin(filename string) {
	if latestImageCid == "" {
		// no image yet
		return
	}

	processInference(latestImageCid)
}

func processImage(filename string) {

	// like /tmp/bacalhau-streaming-cid2357033235/webcam-01/QmcPYjD8R6nmY3pmDDeAPcBNeBRsFjortYsH4sGNaUz7ov/image.jpeg

	shrapnel := strings.Split(filename, "/")
	cid := shrapnel[len(shrapnel)-1]

	postToSlack(fmt.Sprintf(
		"received an image: http://212.82.90.194:9009/%s/image.jpeg", cid,
	))

	latestImageCid = filename
}

func main() {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	postToSlack("🤖 On-prem demo booted 🤖")

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Create) {

					// directory creation trigger

					// cids are gonna show up in /ap_connections/:cid or /webcam_images/:cid

					if strings.HasPrefix(event.Name, "/ap_connections") {
						processAPJoin(event.Name)
					} else if strings.HasPrefix(event.Name, "/webcam_images") {
						processImage(event.Name)
					}

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/ap_connections")
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Add("/webcam_images")
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
