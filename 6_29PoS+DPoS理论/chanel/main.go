package main

import (
	"fmt"
	"time"
)

//线程传参
var data = make(chan int)


func main() {

	go func() {

		for {
			msg:=<-data
			fmt.Println(msg)
			break

		}


	} ()

	go func() {
		time.Sleep(2*time.Second)
		data<-100
	}()

	for   {

	}


}
