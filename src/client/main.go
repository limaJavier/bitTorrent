package main

import (
	"bittorrent/client/peer"
	"bittorrent/common"
	"bittorrent/torrent"
	"fmt"
	"net"
	"os"
	"sync"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("ERROR: expecting 3 arguments: torrent file, download directory and ip")
		os.Exit(1)
		return
	}

	torrentFileName := os.Args[1]
	downloadDirectory := os.Args[2]
	ip := os.Args[3]

	_torrent, err := torrent.ParseTorrentFile(torrentFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	peerId := common.GenerateRandomString(20)
	listener, err := net.Listen("tcp", ip+":")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	_peer, err := peer.New(peerId, listener, _torrent, downloadDirectory, false)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	_peer.Torrent(&wg)
	wg.Wait()
}
