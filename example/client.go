package main

import (
	"fmt"
	"net/http"
	"time"
)

func sendHeartbeat(name string) {

	for {
		_, err := http.Get(fmt.Sprintf("http://localhost:8001/heartbeat?name=%s", name))
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(3 * time.Second)
	}

}

func main() {

	go sendHeartbeat("test1")
	go sendHeartbeat("test2")
	fmt.Println("start")
	select {}
}
