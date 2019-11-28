# About

Simple example of how to broadcast messages two one or multiple subscriber (channels) using Go `channels`.

## Build

```bash
git clone github.com/embano1/broadcaster
cd broadcaster
export GO111MODULE=on # for versions before Go v1.13, not needed if you use $GOPATH
go build -o broadcaster main.go
```

## Run

`main.go` is constructed in a way so that it will immediately add two `channels` to the broadcaster during broadcaster construction. You will see integers being sent to both channels. After some delay, a third channel will be added. This channel blocks for some time, which is logged to standard output. After the blocking delay, messages will also be sent to channel 3.

> **Note:** The example currently uses an array (map with UUID to be implemented) so any log line with a channel ID should be read as starting from index 0, i.e. `outChan3` becomes `channel 2` in a log statement.

## Todo

- [ ] Use a map
- [ ] Add option to remove a channel