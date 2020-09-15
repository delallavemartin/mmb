# imb 0.3
Incremental Message Broker

ab -n 5000 -c 20 -p code/golang-projects/mmb/SampleTextFile_1000kb.tsv http://localhost:8080/notify > code/golang-projects/mmb/incremental-mb/0.2/benchmarks/5000_20_1k_channel.txt

2020/09/15 15:15:37 ERROR - send request to port:  http://localhost:9999/notify  FAILED
2020/09/15 15:15:37 Post "http://localhost:9999/notify": dial tcp: lookup localhost: device or resource busy

9998 port was proccesed for last, but just a bit

