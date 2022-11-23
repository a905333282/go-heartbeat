package main

import (
	"fmt"
	"github.com/a905333282/go-heartbeat/heartbeat_service"
	"net/http"
	"time"
)

var heartbeatService *heartbeat_service.HeartbeatService

func Heartbeat(w http.ResponseWriter, req *http.Request) {

	name := req.URL.Query().Get("name")

	if !heartbeatService.Contains(name) {
		heartbeatService.AddPulse(name)
	}

	heartbeatService.GetAllPulses()

	err := heartbeatService.ResetPulse(name)
	if err != nil {
		return
	}

	fmt.Printf("Get heartbeat from %s", name)
}

func StatesMonitor() {
	for {
		fmt.Println(heartbeatService.GetAllPulsesStatus())
		time.Sleep(3 * time.Second)
	}
}

func main() {

	heartbeatService = heartbeat_service.NewHeartbeatService(heartbeat_service.SetDuration(3 * time.Second))

	go StatesMonitor()

	http.HandleFunc("/heartbeat", Heartbeat)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		return
	}

}
