# GO PPROF Example

Simple example of Gorilla MUX running a simple web server and a Pprof Server.

PProf Server has all endpoints enabled:

* "/"
* "/cmdline"
* "/symbol"
* "/trace"
* "/profile"

https://stackoverflow.com/questions/19591065/profiling-go-web-application-built-with-gorillas-mux-with-net-http-pprof
	
* "/goroutine"
* "/heap"
* "/threadcreate"
* "/block"

## Performance client

I used Apache Benchmark. Some examples in perf.sh

## Dockerfile



