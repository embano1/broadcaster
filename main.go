package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/embano1/broadcaster/broadcaster"
)

func main() {
	outCh1 := make(chan interface{})
	outCh2 := make(chan interface{})
	inCh := make(chan interface{})
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	b, err := broadcaster.New(inCh, &outCh1, &outCh2)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-time.After(30 * time.Second)
		cancel()
	}()

	go func() {
		<-sigCh
		log.Printf("got signal, shutting down")
		cancel()
	}()

	go func() {
		for i := 0; i < 100; i++ {
			inCh <- i
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for v := range outCh1 {
			log.Printf("outCh1 received value: %v", v)
		}
	}()

	go func() {
		for v := range outCh2 {
			log.Printf("outCh2 received value: %v", v)
		}
	}()

	go func() {
		<-time.After(5 * time.Second)
		outCh3 := make(chan interface{})
		b.Add(&outCh3)
		<-time.After(3 * time.Second)
		for v := range outCh3 {
			log.Printf("outCh3 received value: %v", v)
		}
	}()

	err = b.Run(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("done")
}
