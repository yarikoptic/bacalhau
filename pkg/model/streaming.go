package model

type StreamingResult struct {
	// one of these two will be defined
	// LocalPath for images on the filesystem - TODO: it would be nicer if this
	// was the path inside the container, but for now it's the path on the host
	LocalPath string `json:"LocalPath"`
	// InlineData for things like log lines
	InlineData string `json:"InlineData"`
	// the Channel to broadcast the streaming result on
	// this is a logical channel so we can use the same
	// gossip sub topic otherwise we have to manually connect
	Channel string `json:"Channel"`
	// TODO: consider adding tags, e.g. so that webcam can tag its location or
	// identity for the data stream
}
