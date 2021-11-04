package fcfs

import (
	"fmt"
	"math/rand"
	"time"
)

func searchByBing() interface{} {
	return "searchByBing"
}
func searchByGoogle() interface{} {
	return "searchByGoogle"
}

func FirstComeFirstServe() {
	resp := make(chan interface{})
	rand.Seed(time.Now().UnixNano())
	go func() {
		time.Sleep(time.Duration(rand.Intn(1000)))
		resp <- searchByBing()
	}()
	go func() {
		time.Sleep(time.Duration(rand.Intn(1000)))
		resp <- searchByGoogle()
	}()

	result, _ := (<-resp).(string)
	fmt.Println(result)
}
