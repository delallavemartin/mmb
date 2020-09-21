# imb 0.0 - Incremental Message Broker

# SPEC
  1. HTTP Server Up & Running.
  2. Notify endpoint: Will process each event sent by the *Publisher*.
  3. Handler that dispatchs message from *Broker* to *Subscriber*. (It's Supports only one hardcoded subscriber)

# HOW TO

**1. HTTP Server Up and Running**
```go
func main() {
	...
	// Each request its mapped to one lightweight thread trough go routines.
	http.ListenAndServe(":8080", nil)

}
```

**2. Notify endpoint**

```go
func main() {
	// When server receives a notification, the msg will be published to his subscribers/consumers
	http.HandleFunc("/notify", publisherHandler)
	...
}
```

**3. Handler**

```go
func publisherHandler(w http.ResponseWriter, r *http.Request) {
	statusCode := postMsg("http://localhost:9999/notify", "text/plain", r.Body)
	w.WriteHeader(statusCode)
}
```

# Benchmark

ab -n 5000 -c 20 -p code/golang-projects/mmb/SampleTextFile_1000kb.tsv http://localhost:8080/notify > code/golang-projects/mmb/incremental-mb/0.1/benchmarks/5000_20_1k.txt
