# Execution Tracing in Go

This repo showcases how to use the `pprof` and the `trace` tools, built into the go binary, in order to take informed 
decision on how to improve an application.

## How to use

Uncomment the following lines inside the `main.go` file in order to start profiling or tracing:

```go
package main

func main() {
    // Uncomment one at a time

	// pprof.StartCPUProfile(os.Stdout)
	// defer pprof.StopCPUProfile()

	// trace.Start(os.Stdout)
	// defer trace.Stop() 
	
	...
}
```

After doing that rebuild the binary and execute the following commands: 

```bash
$ go build
$ time ./execution-tracing-in-go > t.out # for tracing
$ go tool trace t.out
```
