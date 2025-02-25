package main

import (
	"bittorrent/dht/library"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func AddNode(database *library.DataBaseInMemory, barrier *sync.WaitGroup) {
	randomId := library.GenerateRandomBinaryId()
	fmt.Println("Adding Node ", randomId)
	iString := strconv.Itoa(int(randomId))
	var server = library.NewServerInMemory(database, "Server"+iString)
	var client = library.NewClientInMemory(database, "Client"+iString)
	var monitor = library.NewMonitorHand[library.InMemoryContact]("Monitor" + iString)
	node := library.NewBruteChord[library.InMemoryContact](server, client, monitor, randomId)
	database.AddNode(node, server, client)
	barrier.Add(1)
	go func() {
		node.BeginWorking()
		defer barrier.Done()
	}()
}

func RemoveNode(database *library.DataBaseInMemory, barrier *sync.WaitGroup) {
	if len(database.GetNodes()) > 0 {
		for _, node := range database.GetNodes() {
			barrier.Add(1)
			go func() {
				fmt.Println("Removing Node with ID = ", node.GetId())
				database.RemoveNode(node)
				defer barrier.Done()
			}()
			break
		}
	}
}

func PutData(database *library.DataBaseInMemory) {
	if len(database.GetNodes()) > 0 {
		for _, node := range database.GetNodes() {
			key := rand.Int() % (1 << library.NumberBits)
			val := []byte{byte(key)}
			fmt.Printf("Going to put key %v with data %v using query node = %v \n", key, val, node.GetId())
			node.Put(library.ChordHash(key), val)
			break
		}
	}
}

func ScenarioEasy() (*library.DataBaseInMemory, *sync.WaitGroup) {
	library.SetLogDirectoryPath("BasicScenario")
	var database = *library.NewDataBaseInMemory()
	var barrier = sync.WaitGroup{}
	fmt.Println("Nodes are being added and removed randomly every once a while")
	barrier.Add(1)
	go func() {
		defer barrier.Done()
		for {
			time.Sleep(1 * time.Second)
			if rand.Float32() <= 0.3 {
				AddNode(&database, &barrier)
			}
			if rand.Float32() <= 0.1 {
				RemoveNode(&database, &barrier)
			}
		}
	}()
	return &database, &barrier
}

func ScenarioMedium() (*library.DataBaseInMemory, *sync.WaitGroup) {
	library.SetLogDirectoryPath("MediumScenario")
	var database = *library.NewDataBaseInMemory()
	var barrier = sync.WaitGroup{}
	barrier.Add(1)
	go func() {
		defer barrier.Done()
		for {
			time.Sleep(3 * time.Second)
			if rand.Float32() <= 0.4 {
				AddNode(&database, &barrier)
			}
			if rand.Float32() <= 0.1 {
				RemoveNode(&database, &barrier)
			}
			if rand.Float32() <= 0.7 {
				PutData(&database)
			}
		}
	}()
	return &database, &barrier
}

func PutScenario() (*library.DataBaseInMemory, *sync.WaitGroup) {
	N := 10
	library.SetLogDirectoryPath("PutScenario")
	var database = *library.NewDataBaseInMemory()
	var barrier = sync.WaitGroup{}
	var toPut = make(map[library.ChordHash][]byte)
	fmt.Printf("Going to Add N = %v  Nodes", N)
	for i := 0; i < N; i++ {
		AddNode(&database, &barrier)
	}
	for i := 0; i < 50; i++ {
		toPut[library.ChordHash(i)] = []byte{byte(i)}
	}
	time.Sleep(2 * time.Second)
	nodes := database.GetNodes()
	println(len(nodes))
	for key, value := range toPut {
		for _, node := range nodes {
			go node.Put(key, value)
		}
	}
	return &database, &barrier
}

func main() {
	database, barrier := ScenarioMedium()
	library.StartGUI(database, barrier)
}
