package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"os/exec"
	"time"

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
	resp, err := http.Post(os.Getenv("SLACK_WEBHOOK_URL"), "application/json", strings.NewReader(payload))
	if err != nil {
		log.Printf("error posting to slack: %s", err)
		return err
	}
	defer resp.Body.Close()
	log.Printf("posted to slack: %s", text)
	// read response body
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		return err
	}
	log.Printf("slack response: %s", string(bs))
	return nil
}

func notifySlack(labels []string) {
	// make an http request to slack webhook
	postToSlack(fmt.Sprintf("detected %+v!", labels))
}

func processInference(latestImageCid string) {
	// make directory /outputs/:latestImageCid
	_, err := os.Stat("/outputs/" + latestImageCid)
	if err != nil {
		log.Printf("error checking for directory: %s", err)
		return
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll("/outputs/"+latestImageCid, 0755)
	} else if err != nil {
		log.Printf("error creating directory: %s", err)
		return
	}

	postToSlack(fmt.Sprintf("RUNNING INFERENCE 🤔..."))
	// run python detect.py --weights /weights/yolov5s.pt --source /webcam_images/QmeEjqtVU2dZsPpUn1r8cyNXq8ptTwzKKnzv57oUx5Ru7R/ --project /outputs/QmeEjqtVU2dZsPpUn1r8cyNXq8ptTwzKKnzv57oUx5Ru7R
	log.Printf("running inference on %s", latestImageCid)
	log.Printf("about to run python detect.py --weights /weights/yolov5s.pt --source /webcam_images/%s/ --project /outputs/%s", latestImageCid, latestImageCid)
	output, err := exec.Command(
		"python", "detect.py", "--weights", "/weights/yolov5s.pt",
		"--source", "/webcam_images/"+latestImageCid+"/",
		"--project", "/outputs/"+latestImageCid,
	).CombinedOutput()
	log.Printf("output: %s", output)
	if err != nil {
		log.Printf("error running inference: %s", err)
		return
	}

	// do AI inference to get labels
	postToSlack(fmt.Sprintf(
		"INFERENCE: http://mind.lukemarsden.net:9010/%s/image.jpeg", latestImageCid,
	))
	postToSlack(fmt.Sprintf(
		"OUTPUT:\n```\n%s\n```", output,
	))
}

func processAPJoin(filename string) {

	// ugh, the directory is created before the file inside it exists
	time.Sleep(100 * time.Millisecond)

	log.Printf("processAPJoin called with %s", filename)
	// read filename into string
	f, err := os.Open(filename + "/output.txt")
	if err != nil {
		log.Printf("error opening file: %s", err)
		return
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("error reading file: %s", err)
		return
	}

	if latestImageCid == "" {
		// no image yet
		log.Printf("no image yet, skipping inference")
		return
	}

	postToSlack(fmt.Sprintf(
		"ACCESS POINT CONNECTION DETECTED: %s",
		string(bs),
	))

	postToSlack(fmt.Sprintf(
		"LATEST IMAGE: http://mind.lukemarsden.net:9009/%s/image.jpeg", latestImageCid,
	))

	processInference(latestImageCid)
}

func processImage(filename string) {

	// like /tmp/bacalhau-streaming-cid2357033235/webcam-01/QmcPYjD8R6nmY3pmDDeAPcBNeBRsFjortYsH4sGNaUz7ov/image.jpeg

	shrapnel := strings.Split(filename, "/")
	cid := shrapnel[len(shrapnel)-1]

	latestImageCid = cid
}

func main() {
	log.Printf("Starting up.")
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
