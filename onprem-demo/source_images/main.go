package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	imageSourceDir := os.Getenv("ONPREM_DEMO_IMAGE_SOURCE_DIR")
	if imageSourceDir == "" {
		msg := "Please specify ONPREM_DEMO_IMAGE_SOURCE_DIR as an absolute path to where " +
			"your webcam image software is writing jpeg files, e.g. \n\n\t" +
			"while true; do streamer -c /dev/video0 -b 16 -o /path/to/cam01-$(date +%s).jpeg; sleep 1; done\n"
		fmt.Println(msg)
		os.Exit(1)
	}
	// Add a path.
	err = watcher.Add(imageSourceDir)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
