package main

import (
	"fmt"
	"github.com/a905333282/go-heartbeat/heartbeat_service"
	"time"
)

func main() {
	heartbeatService := heartbeat_service.NewHeartbeatService()
	heartbeatService.AddPulse("test")
	heartbeatService.AddPulse("test1")
	fmt.Println(heartbeatService.GetAllPulsesStatus())
	go func() {
		for {
			timer := time.After(time.Second * 1)
			<-timer
			err := heartbeatService.ResetPulse("test")
			if err != nil {
				return
			}
		}
	}()
	go func() {
		for i := 0; i < 5; i++ {
			timer := time.After(time.Second * 1)
			<-timer
			err := heartbeatService.ResetPulse("test1")
			if err != nil {
				return
			}
		}
	}()
	timer := time.After(time.Second * 10)
	<-timer
	fmt.Println(heartbeatService.GetAllPulsesStatus())
	err := heartbeatService.RemovePulse("test1")
	if err != nil {
		fmt.Println(err)
	}
	select {}
	//fmt.Println("start")
	//ch := make(chan interface{}, 10)
	//go heartbeat(ch)
	//MyPulse := NewPulse(ch, 3*time.Second)
	//MyPulse.Start()
	//close(ch)
	//select {}
	//fmt.Println("start")
	//
	//ch := make(chan interface{}, 10) // 可读可写且带缓存区
	//go heartbeat(ch)                 // 模拟客户端
	//go heartbeatHandler(ch)
	//
	//select {}

}

func heartbeatHandler(ch chan interface{}) {
	for {
		timer := time.NewTimer(3 * time.Second)
		select {
		case str := <-ch:
			fmt.Println("receive str", str)
			timer.Reset(3 * time.Second)
		case <-timer.C:
			fmt.Println("timeout!!")
		}
	}
}

func heartbeat(ch chan interface{}) {
	for i := 0; i < 5; i++ {
		timer := time.After(time.Second * 1)
		<-timer
		ch <- "heartbeat happened!"
	}
}
