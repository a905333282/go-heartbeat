package main

import (
	"fmt"
	"github.com/a905333282/go-heartbeat/heartbeat_service"
	"github.com/a905333282/go-heartbeat/hook"
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

	//fmt.Printf("Get heartbeat from %s", name)
}

func StatesMonitor() {
	for {
		fmt.Println(heartbeatService.GetAllPulsesStatus())
		heartbeatService.RemovePulse("test1")
		time.Sleep(1 * time.Second)
	}
}

type MyHook struct {
	hook.DefaultHooks
}

func (mh *MyHook) OnAdd(name string) {
	fmt.Println("OnAdd:", name)
}
func (mh *MyHook) OnDeath(name string) {
	fmt.Println("OnDeath:", name)
}
func (mh *MyHook) OnReset(name string) {
	fmt.Println("OnReset:", name)
}
func (mh *MyHook) OnRenewal(name string) {
	fmt.Println("OnRenewal:", name)
}
func (mh *MyHook) OnRemove(name string) {
	fmt.Println("OnRemove:", name)
}

func main() {

	heartbeatService = heartbeat_service.NewHeartbeatService(heartbeat_service.SetHooks(&MyHook{}))

	go StatesMonitor()

	http.HandleFunc("/heartbeat", Heartbeat)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		return
	}

}
