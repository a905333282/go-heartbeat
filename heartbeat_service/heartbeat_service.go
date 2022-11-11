package heartbeat_service

import (
	"errors"
	"fmt"
	"github.com/a905333282/go-heartbeat/pulse"
	"time"
)

type HeartbeatService struct {
	Pulses   map[string]*pulse.Pulse
	Channels map[string]chan interface{}
}

// AddPulse
// Add a new Pulse with a unique name in string
func (h *HeartbeatService) AddPulse(name string) {
	ch := make(chan interface{}, 10)
	pulse := pulse.NewPulse(name, ch, 3*time.Second)

	h.Channels[name] = ch
	h.Pulses[name] = pulse

	go pulse.Start()
}

// RemovePulse
// Remove a Pulse
func (h *HeartbeatService) RemovePulse(name string) error {
	if _, ok := h.Channels[name]; ok {
		h.Pulses[name].Stop()
		delete(h.Pulses, name)
		delete(h.Channels, name)
		return nil
	}
	return errors.New(fmt.Sprintf("the Pulse '%s' didn't register", name))
}

// ResetPulse
// Call this function with a pulse name when get a heartbeat
func (h *HeartbeatService) ResetPulse(name string) error {

	if _, ok := h.Channels[name]; ok {
		h.Channels[name] <- nil

		return nil
	}
	return errors.New(fmt.Sprintf("the Pulse '%s' didn't register", name))
}

// NumberOfPulse
// Get the number of pulse
func (h *HeartbeatService) NumberOfPulse() int {
	return len(h.Pulses)
}

func (h *HeartbeatService) GetAllPulses() []string {
	index := 0
	pulses := make([]string, len(h.Pulses))
	for key := range h.Pulses {
		pulses[index] = key
		index++
	}
	return pulses
}

// GetAllPulsesStatus
// 得到所有Pulse的状态，ALIVE或者DEAD, 0:ALIVE, 1:DEAD
func (h *HeartbeatService) GetAllPulsesStatus() map[string]pulse.StateType {
	pulses := make(map[string]pulse.StateType, len(h.Pulses))
	for key := range h.Pulses {
		pulses[key] = h.Pulses[key].GetState()
	}
	return pulses
}

func NewHeartbeatService() *HeartbeatService {
	return &HeartbeatService{
		Pulses:   make(map[string]*pulse.Pulse),
		Channels: make(map[string]chan interface{}),
	}
}
