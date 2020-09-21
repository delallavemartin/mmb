# imb 0.2 - Incremental Message Broker

# SPEC

1. Performance still not good enough

# HOW TO: Step 1 - CHANNELS

```go
func publish(ch chan Request) {
	// this loop receives values from the channel repeatedly until it is closed
	for request := range ch {
		// go routines added to improve request per second performance.
		go postMsg(request.Url, request.ContentType, request.Reader)
	}

}

func publisherHandler(w http.ResponseWriter, r *http.Request) {
	// Reader converted to string to create one Reader per POST.
	// Since is requested by postMsg function firm
	body := readerToString(r.Body)

	// Channel
	ch := make(chan Request)

	go publish(ch)

	// Sends requests to the channel in order to proccess them.
	ch <- Request{"http://localhost:9995/notify", "text/plain", strings.NewReader(body)}
	ch <- Request{"http://localhost:9996/notify", "text/plain", strings.NewReader(body)}
	ch <- Request{"http://localhost:9997/notify", "text/plain", strings.NewReader(body)}
	ch <- Request{"http://localhost:9998/notify", "text/plain", strings.NewReader(body)}
	ch <- Request{"http://localhost:9999/notify", "text/plain", strings.NewReader(body)}

}
```

# Benchmark: Step 1

ab -n 5000 -c 20 -p code/golang-projects/mmb/SampleTextFile_1000kb.tsv http://localhost:8080/notify > code/golang-projects/mmb/incremental-mb/0.2/benchmarks/5000_20_1k_channel.txt

2020/08/11 15:21:08 ERROR - send request to port:  http://localhost:9997/notify  FAILED
2020/08/11 15:21:08 Post "http://localhost:9997/notify": dial tcp: lookup localhost: device or resource busy

2020/08/11 15:21:08 http: Accept error: accept tcp [::]:8080: accept4: too many open files; retrying in 10ms

9995 port was proccesed for last, but just a bit, the load was more distributed


# HOW TO: Step 2 - BUFFERED CHANNELS

```go
func publish(ch chan Request) {
	// this loop receives values from the channel repeatedly until it is closed
	for request := range ch {
		// go routines added to improve request per second performance.
		go postMsg(request.Url, request.ContentType, request.Reader)
	}

}

func publisherHandler(w http.ResponseWriter, r *http.Request) {
	// Reader converted to string to create one Reader per POST.
	// Since is requested by postMsg function firm
	body := readerToString(r.Body)

	// Buffered Channel
	ch := make(chan Request, 5)

	go publish(ch)

	// Sends requests to the channel in order to proccess them.
	ch <- Request{"http://localhost:9995/notify", "text/plain", strings.NewReader(body)}
	ch <- Request{"http://localhost:9996/notify", "text/plain", strings.NewReader(body)}
	ch <- Request{"http://localhost:9997/notify", "text/plain", strings.NewReader(body)}
	ch <- Request{"http://localhost:9998/notify", "text/plain", strings.NewReader(body)}
	ch <- Request{"http://localhost:9999/notify", "text/plain", strings.NewReader(body)}

}
```
# Benchmark: Step 1

ab -n 5000 -c 20 -p code/golang-projects/mmb/SampleTextFile_1000kb.tsv http://localhost:8080/notify > code/golang-projects/mmb/incremental-mb/0.2/benchmarks/5000_20_1k_bf_channel.txt

2020/08/11 15:21:08 ERROR - send request to port:  http://localhost:9997/notify  FAILED
2020/08/11 15:21:08 Post "http://localhost:9997/notify": dial tcp: lookup localhost: device or resource busy

2020/08/11 15:21:08 http: Accept error: accept tcp [::]:8080: accept4: too many open files; retrying in 10ms

9995 port was proccesed for last, same values as it doesnt have channel.



