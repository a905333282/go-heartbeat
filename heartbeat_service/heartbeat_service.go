package heartbeat_service

import (
	"errors"
	"fmt"
	"github.com/a905333282/go-heartbeat/hook"
	"github.com/a905333282/go-heartbeat/pulse"
	"time"
)

type HeartbeatService struct {
	Pulses       map[string]*pulse.Pulse
	Channels     map[string]chan interface{}
	Duration     time.Duration // 多少秒收不到heartbeat就改变pulse状态
	UniformHooks hook.Hooks
}

// AddPulse
// Add a new Pulse with a unique name in string
func (h *HeartbeatService) AddPulse(name string) {
	ch := make(chan interface{}, 10)
	newPulse := pulse.NewPulse(name, ch, h.Duration, h.UniformHooks)
	h.Channels[name] = ch
	h.Pulses[name] = newPulse
	// 调用hook函数
	go h.UniformHooks.OnAdd(name)

	go newPulse.Start()
}

// RemovePulse
// Remove a Pulse
func (h *HeartbeatService) RemovePulse(name string) error {
	if _, ok := h.Channels[name]; ok {
		h.Pulses[name].Stop()
		delete(h.Pulses, name)
		delete(h.Channels, name)
		// 调用hook函数
		go h.UniformHooks.OnRemove(name)
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

func (h *HeartbeatService) Contains(name string) bool {

	if _, ok := h.Pulses[name]; ok {
		return true
	}
	return false
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

// HeartbeatServiceOption 定义一个所有默认配置的结构体
type HeartbeatServiceOption struct {
	Duration     time.Duration // 多少秒收不到heartbeat就改变pulse状态
	UniformHooks hook.Hooks
}

// Option 定义一个接口，实现该接口需要实现apply函数
type Option interface {
	apply(*HeartbeatServiceOption)
}

// funcOption 实现了上述的接口Option
type funcOption struct {
	f func(*HeartbeatServiceOption)
}

func (fdo *funcOption) apply(do *HeartbeatServiceOption) {
	fdo.f(do)
}

// 新建一个funcOption， 传入一个函数作为参数
func newFuncOption(f func(*HeartbeatServiceOption)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func SetDuration(duration time.Duration) Option {
	return newFuncOption(func(o *HeartbeatServiceOption) {
		o.Duration = duration
	})
}

func SetHooks(uniformHooks hook.Hooks) Option {
	return newFuncOption(func(o *HeartbeatServiceOption) {
		o.UniformHooks = uniformHooks
	})
}

func DefaultHeartbeatServiceOption() HeartbeatServiceOption {
	return HeartbeatServiceOption{
		Duration:     9 * time.Second,
		UniformHooks: &hook.DefaultHooks{},
	}
}

func NewHeartbeatService(opts ...Option) *HeartbeatService {

	defaultOpts := DefaultHeartbeatServiceOption()
	for _, opt := range opts {
		opt.apply(&defaultOpts)
	}

	return &HeartbeatService{
		Pulses:       make(map[string]*pulse.Pulse),
		Channels:     make(map[string]chan interface{}),
		Duration:     defaultOpts.Duration,
		UniformHooks: defaultOpts.UniformHooks,
	}
}
