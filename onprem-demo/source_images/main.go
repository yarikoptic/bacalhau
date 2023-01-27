package main

import (
	"fmt"
	"io"
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

	imageSinkDir := os.Getenv("ONPREM_DEMO_IMAGE_SINK_DIR")
	if imageSinkDir == "" {
		msg := "Please specify ONPREM_DEMO_IMAGE_SINK_DIR as an absolute path to where " +
			"you have configured your Bacalhau job's output directory, e.g. /output"
		fmt.Println(msg)
		os.Exit(1)

	}
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
					// XXX what if the whole file write isn't finished yet
					err, _ := copy(event.Name, imageSinkDir+"/"+"image.jpg")
					if err != nil {
						log.Printf("error copying file: %s", err)
					}
					// TODO: http request
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

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
