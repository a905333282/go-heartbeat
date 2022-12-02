package pulse

import (
	"github.com/a905333282/go-heartbeat/hook"
	"time"
)

type StateType int32

const (
	ALIVE StateType = 0
	DEAD  StateType = 1
)

type Pulse struct {
	Channel  chan interface{}
	State    StateType
	Timer    *time.Timer
	Duration time.Duration
	Name     string
	StopCh   chan interface{}
	Hooks    hook.Hooks
}

func (p *Pulse) Start() {
	for {
		select {

		case _ = <-p.StopCh:
			goto END

		case _ = <-p.Channel:

			go p.Hooks.OnReset(p.Name)

			if p.State == DEAD {
				go p.Hooks.OnRenewal(p.Name)
			}

			p.State = ALIVE
			p.Timer.Reset(p.Duration)

		case <-p.Timer.C:

			if p.State == ALIVE {
				go p.Hooks.OnDeath(p.Name)
			}

			p.State = DEAD
			p.Timer.Reset(p.Duration)
		}
	}
END:
}

func (p *Pulse) Stop() {
	p.StopCh <- nil
	close(p.StopCh)
	close(p.Channel)
}

func (p *Pulse) GetState() StateType {
	return p.State
}

func NewPulse(name string, channel chan interface{}, duration time.Duration, hooks hook.Hooks) *Pulse {
	return &Pulse{
		Channel:  channel,
		State:    ALIVE,
		Timer:    time.NewTimer(duration),
		Duration: duration,
		Name:     name,
		StopCh:   make(chan interface{}, 1),
		Hooks:    hooks,
	}
}
