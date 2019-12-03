# About

Simple example of how to broadcast messages on one input channel to zero or multiple subscribers using Go `channels`.

## Build

```bash
git clone github.com/embano1/broadcaster
cd broadcaster
export GO111MODULE=on # for versions before Go v1.13, not needed if you use $GOPATH
go build -o broadcaster main.go
```

## Run

`main.go` is constructed in a way so that it will immediately add two `subscribers` to the broadcaster during broadcaster construction. 

You will see integers being sent to both channels. After some delay, a third subscriber will be added. Its channel blocks for some time, which is logged to standard output. After the blocking delay, messages will also be sent to subscriber 3. Meanwhile subscriber 2 will unsubscribe from the broadcaster and then safely close its channel used. After a timer the context is cancelled via `cancel()` and the broadcaster shuts down.

It is possible to use buffered channels to reduce the chance of missing events when a subscriber would be blocked.

> **Note:** Due to the indeterministic runtime behavior you might see the first subscribers also being blocked on the first message(s) sent. This is due to the way the mutex is used and how the runtime spawns and schedules the individual goroutines.
