# imb 0.1 - Incremental Message Broker

# SPEC

1. Supports five different _Subscribers_

# HOW TO: Step 1

```go
func publisherHandler(w http.ResponseWriter, r *http.Request) {
	// Reader converted to string to create one Reader per POST.
	// Since is requested by postMsg function firm
	body := readerToString(r.Body)
	postMsg("http://localhost:9995/notify", "text/plain", strings.NewReader(body))
	postMsg("http://localhost:9996/notify", "text/plain", strings.NewReader(body))
	postMsg("http://localhost:9997/notify", "text/plain", strings.NewReader(body))
	postMsg("http://localhost:9998/notify", "text/plain", strings.NewReader(body))
	postMsg("http://localhost:9999/notify", "text/plain", strings.NewReader(body))
}
```

# Benchmark: Step 1

ab -n 5000 -c 20 -p code/golang-projects/mmb/SampleTextFile_1000kb.tsv http://localhost:8080/notify > code/golang-projects/mmb/incremental-mb/0.1/benchmarks/5000_20_1k.txt

too slow

# HOW TO: Step 2

```go
func publisherHandler(w http.ResponseWriter, r *http.Request) {
	// Reader converted to string to create one Reader per POST.
	// Since is requested by postMsg function firm
	body := readerToString(r.Body)
	go postMsg("http://localhost:9995/notify", "text/plain", strings.NewReader(body))
	go postMsg("http://localhost:9996/notify", "text/plain", strings.NewReader(body))
	go postMsg("http://localhost:9997/notify", "text/plain", strings.NewReader(body))
	go postMsg("http://localhost:9998/notify", "text/plain", strings.NewReader(body))
	go postMsg("http://localhost:9999/notify", "text/plain", strings.NewReader(body))
}

```
# Benchmark: Step 1

ab -n 5000 -c 20 -p code/golang-projects/mmb/SampleTextFile_1000kb.tsv http://localhost:8080/notify > code/golang-projects/mmb/incremental-mb/0.1/benchmarks/5000_20_1k_go_routines.txt

better by far, but...

2020/08/11 15:21:08 ERROR - send request to port:  http://localhost:9997/notify  FAILED
2020/08/11 15:21:08 Post "http://localhost:9997/notify": dial tcp: lookup localhost: device or resource busy

2020/08/11 15:21:08 http: Accept error: accept tcp [::]:8080: accept4: too many open files; retrying in 10ms

9995 port was proccesed for last
