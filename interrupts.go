package gointhandler

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type AsyncSignalHandlerFunc func(sig os.Signal)

var _stop bool = false

var StopPointer *bool = &_stop
var StopChannel = make(chan struct{}, 10) // ..should be enough? anyways will just close it
var MaxTimeToWaitForCleanup = time.Second * 5
var AsyncSignalHandler AsyncSignalHandlerFunc = nil

var interruptChannel = make(chan os.Signal, 1)

func InterruptTheApp() {
	interruptChannel <- os.Interrupt
	time.Sleep(time.Millisecond * 20)
}

// convenience
func IsTimeToStop() bool {
	return *StopPointer
}

func TakeCareOfInterrupts(ignoreSIGALRM bool) {
	c := interruptChannel
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGALRM)

	go func() {
		defer func() { recover() }()
		for sig := range c {
			log.Println("got signal:", sig.String())
			if sig == syscall.SIGHUP {
				if AsyncSignalHandler != nil {
					go AsyncSignalHandler(sig)
				}
				continue
			}
			if sig == syscall.SIGALRM && !ignoreSIGALRM {
				log.Fatalln("need to exit immediately")
				//log.Flush() -- no flush, no wait!
			}

			if !(*StopPointer) {
				*StopPointer = true
				go func() {
					close(StopChannel)
					time.Sleep(MaxTimeToWaitForCleanup)
					log.Fatalln("forced shutdown after waiting for", MaxTimeToWaitForCleanup.Seconds(), "seconds")
					os.Exit(0)
				}()
			}
		}
		log.Println("this code should never execute.. however who knows.")
		time.Sleep(MaxTimeToWaitForCleanup)
		log.Fatalln("forced shutdown after waiting for", MaxTimeToWaitForCleanup.Seconds(), "seconds")
		os.Exit(0)
	}()
}

var CancelError = errors.New("canceled")

// returns true if was canceled
func Sleep(duration time.Duration, cancelChannel ...chan struct{}) error {
	var cc chan struct{} = make(chan struct{})
	if len(cancelChannel) > 0 {
		cc = cancelChannel[0]
	}
	select {
	case <-time.After(duration):
		return nil
	case <-StopChannel:
		return errors.New("time to shutdown")
	case <-cc:
		return CancelError
	}
}
