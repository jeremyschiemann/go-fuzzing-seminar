package main

import (
	"fmt"
	"sync"
)

func sendData(ch chan int, wg *sync.WaitGroup, data int) {
	defer wg.Done()
	ch <- data
}

func main() {
	ch := make(chan int, 2)
	var wg sync.WaitGroup

	// Sending data concurrently
	wg.Add(2)
	go sendData(ch, &wg, 1)
	go sendData(ch, &wg, 2)

	// Closing the channel after sending
	wg.Wait()
	close(ch)

	// Receiving data
	for v := range ch {
		fmt.Println(v)
	}
}
