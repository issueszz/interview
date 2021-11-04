package fcfs

import (
	"math/rand"
	"time"
)

func searchByBing() string {
	return "searchByBing"
}
func searchByGoogle() string {
	return "searchByGoogle"
}

func FirstComeFirstServe() string {
	resp := make(chan string)
	rand.Seed(time.Now().UnixNano())
	go func() {
		time.Sleep(time.Duration(rand.Intn(1000)))
		resp <- searchByBing()
	}()
	go func() {
		time.Sleep(time.Duration(rand.Intn(1000)))
		resp <- searchByGoogle()
	}()
	return <-resp
}
