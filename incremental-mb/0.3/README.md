# imb 0.3 - Incremental Message Broker

# SPEC 
1. Register _10 Subscribers_ programatically
	1. Subscribe Endpoint
	1. Subscribe Handler
	1. Subscriber List
  
# HOW TO: Step 1

**1.1** 

```go
func main() {
	...
	// When server receives a subscription, port will be added to subscribers list
	http.HandleFunc("/subscribe", subscriberHandler)
	...
}
```


**1.2** 

```go
func subscriberHandler(w http.ResponseWriter, r *http.Request) {
	// Read port number
	port_number := readerToString(r.Body)

	// Append new subscribers to subscribers list
	consumers_adresses = append(consumers_adresses, port_number)
	log.Println("INFO - port added: " + port_number)
}
```

**1.3** 

```go
// Global Variables
// SPEC: Only 10 subscribers are supported
var consumers_adresses = make([]string, 0, 10)
```

```go
// Request Handlers
func publisherHandler(w http.ResponseWriter, r *http.Request) {
	// Reader converted to string to create one Reader per POST.
	// Since is requested by postMsg function firm
	body := readerToString(r.Body)

	ch := make(chan Request)

	go publish(ch)

	//Iterate over subscribers list.
	for i := 0; i < len(consumers_adresses); i++ {
		// Send requests to the channel in order to proccess it.
		ch <- Request{"http://localhost:" + consumers_adresses[i] + "/notify", "text/plain", strings.NewReader(body)}
	}
}
```

# Benchmarks: Step 1

ab -n 5000 -c 20 -p code/golang-projects/mmb/SampleTextFile_1000kb.tsv http://localhost:8080/notify > code/golang-projects/mmb/incremental-mb/0.3/benchmarks/5000_20_1k_channel.txt

2020/09/15 15:15:37 ERROR - send request to port:  http://localhost:9999/notify  FAILED
2020/09/15 15:15:37 Post "http://localhost:9999/notify": dial tcp: lookup localhost: device or resource busy

9998 port was proccesed for last, but just a bit

ab -n 5000 -c 20 -p code/golang-projects/mmb/SampleTextFile_1000kb.tsv http://localhost:8080/notify > code/golang-projects/mmb/incremental-mb/0.3/benchmarks/5000_20_1k_added_middle_of_processing.txt

Same behavior as the previous one.
Last consumer added was executed ok.

# HOW TO: Step 2

**1.2** 

```go
func subscriberHandler(w http.ResponseWriter, r *http.Request) {
	// Read port number
	port_number := readerToString(r.Body)

	// Append new subscribers to subscribers list
	subscribers_addresses.add(port_number)
	log.Println("INFO - port added: " + port_number)
}
```

**1.3** 

```go
// Subscriber wrapper for Sync. purposes
type SafeSubscribers struct {
	addresses []string
	mux       sync.Mutex
}

// SafeSubscribers methods
func (s *SafeSubscribers) add(address string) {
	s.mux.Lock()
	s.addresses = append(s.addresses, address)
	s.mux.Unlock()
}

func (s *SafeSubscribers) getByIndex(index int) string {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.addresses[index]
}

func (s *SafeSubscribers) numberOfSubscribers() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return len(s.addresses)
}
```

```go
// Global Variables
// SPEC: Only 10 subscribers are supported
var subscribers_addresses = SafeSubscribers{addresses: make([]string, 0, 10)}
```

```go
// Request Handlers
func publisherHandler(w http.ResponseWriter, r *http.Request) {
	// Reader converted to string to create one Reader per POST.
	// Since is requested by postMsg function firm
	body := readerToString(r.Body)

	ch := make(chan Request)

	go publish(ch)

	//Iterate over subscribers list.
	for i := 0; i < subscribers_addresses.numberOfSubscribers(); i++ {
		// Send request to the channel in order to proccess it.
		ch <- Request{"http://localhost:" + subscribers_addresses.getByIndex(i) + "/notify", "text/plain", strings.NewReader(body)}
	}
}
```

# Benchmarks: Step 2

ab -n 5000 -c 20 -p code/golang-projects/mmb/SampleTextFile_1000kb.tsv http://localhost:8080/notify > code/golang-projects/mmb/incremental-mb/0.3/benchmarks/5000_20_1k_sync.txt

9998 port was proccesed for last, but just a bit


