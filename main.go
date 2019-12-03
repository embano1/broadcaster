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

	b := broadcaster.New(inCh)
	sub1 := b.Subscribe(&outCh1)
	log.Printf("created new subscriber subscriber1 (%s)", sub1)
	sub2 := b.Subscribe(&outCh2)
	log.Printf("created new subscriber subscriber2 (%s)", sub2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// simulate cancel after 30s
	go func() {
		<-time.After(20 * time.Second)
		log.Printf("time to shut down")
		cancel()
	}()

	// signal handling
	go func() {
		<-sigCh
		log.Printf("got signal, shutting down")
		cancel()
	}()

	// simulate input on inCh
	go func() {
		for i := 0; i < 100; i++ {
			inCh <- i
			time.Sleep(time.Second)
		}
	}()

	// spawn subscriber1
	go func() {
		for v := range outCh1 {
			log.Printf("subscriber1 received value: %v", v)
		}
	}()

	// spawn subscriber2
	go func() {
		for v := range outCh2 {
			log.Printf("subscriber2 received value: %v", v)
		}
	}()

	// subscribe and spawn subscriber3 - this one blocks for some time to
	// simulate full channel behavior
	go func() {
		<-time.After(5 * time.Second)
		outCh3 := make(chan interface{})
		log.Printf("creating new subscriber subscriber3")
		b.Subscribe(&outCh3)
		<-time.After(3 * time.Second)
		for v := range outCh3 {
			log.Printf("subscriber3 received value: %v", v)
		}
	}()

	// unsubscribe subscriber2
	go func() {
		<-time.After(7 * time.Second)
		log.Printf("unsubscribing subscriber2 (%s)", sub2)
		b.Unsubscribe(sub2)
		close(outCh2)
	}()

	b.Run(ctx)
	log.Println("done")
}
