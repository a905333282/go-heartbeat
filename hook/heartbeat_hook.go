package hook

type Hooks interface {
	OnAdd(name string)
	OnDeath(name string)
	OnReset(name string)
	OnRenewal(name string)
	OnRemove(name string)
}

type DefaultHooks struct{}

func (dh *DefaultHooks) OnAdd(name string) {
}
func (dh *DefaultHooks) OnDeath(name string) {
}
func (dh *DefaultHooks) OnReset(name string) {
}
func (dh *DefaultHooks) OnRenewal(name string) {
}
func (dh *DefaultHooks) OnRemove(name string) {
}
