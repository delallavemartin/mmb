package subscribers

import (
	"sync"
)

// MODEL DEFINITION
// SubscribersList with Sync.
type SubscribersList struct {
	Addresses []string
	Mux       sync.Mutex
}

// SubscribersList methods
func (self *SubscribersList) Add(address string) {
	self.Mux.Lock()
	self.Addresses = append(self.Addresses, address)
	self.Mux.Unlock()
}

func (self *SubscribersList) NotifySubscribers(notify func(address string)) {
	self.Mux.Lock()
	defer self.Mux.Unlock()
	for i := 0; i < len(self.Addresses); i++ {
		notify(self.Addresses[i])
	}
}
