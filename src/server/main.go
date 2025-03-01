package main

import (
	"bittorrent/dht/library/WithSocket"
	"bittorrent/server/TrackerNode"
	"fmt"
	"time"
)

func main() {
	WithSocket.RegisterStartUp("eth0", "HttpChord", []string{"12345"})
	var tracker1 = TrackerNode.NewHttpTracker("TrackerDocker")
	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("Tracker is running %v:%v \n", tracker1.Ip, tracker1.Port)
	}
}
