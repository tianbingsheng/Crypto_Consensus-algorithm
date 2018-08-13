package main

import (
	"sync"
	"fmt"
	"time"
)

var mux = &sync.Mutex{}

var cnt =0

func Count() {

	mux.Lock()
	cnt++
	time.Sleep(300*time.Millisecond)
	fmt.Println(cnt)
	mux.Unlock()
}


func main() {

	go Count()
	go Count()

	for   {
	}

}
