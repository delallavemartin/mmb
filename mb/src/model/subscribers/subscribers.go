package subscribers

// Package subscribers defines the data structure and methods for managing a list of subscribers.

import (
	"sync"
)

// SubscribersList manages a thread-safe list of subscriber addresses.
type SubscribersList struct {
	Addresses []string   // Addresses stores the list of subscriber addresses.
	Mux       sync.Mutex // Mux provides mutual exclusion for concurrent access to Addresses.
}

// Add appends a new subscriber address to the list in a thread-safe manner.
func (self *SubscribersList) Add(address string) {
	self.Mux.Lock()
	self.Addresses = append(self.Addresses, address)
	self.Mux.Unlock()
}

// NotifySubscribers iterates through the list of subscribers and calls the provided 'notify' function for each address.
// It ensures thread-safe access to the subscriber list during iteration.
func (self *SubscribersList) NotifySubscribers(notify func(address string)) {
	self.Mux.Lock()
	defer self.Mux.Unlock()
	for i := 0; i < len(self.Addresses); i++ {
		notify(self.Addresses[i])
	}
}
