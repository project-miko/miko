package timer

import (
	"sync/atomic"
	"time"

	"github.com/project-miko/miko/tools/errutils"
)

type TimerEvent struct {
	callBack    func()       // callback function
	delay       uint32       // execution interval
	repeatCount uint32       // execution times
	ticker      *time.Ticker // timer
	closeFlag   int32        // close flag
	closeChan   chan int     // close channel
}

// check if it is closed
func (this *TimerEvent) IsClosed() bool {
	return atomic.LoadInt32(&this.closeFlag) == 1
}

// close
func (this *TimerEvent) Close() {
	if atomic.CompareAndSwapInt32(&this.closeFlag, 0, 1) {
		this.ticker.Stop()
		close(this.closeChan)
	}
}

// execute indefinitely
func DoTimer(delay uint32, callback func()) *TimerEvent {
	return Do(delay, 0, callback)
}

// delay processing
func SetTimeOut(delay uint32, callback func()) *TimerEvent {
	return Do(delay, 1, callback)
}

// remove a timer
func Remove(event *TimerEvent) {
	if event == nil {
		return
	}
	event.Close()
}

// time interval, execution times, callback function
func Do(delay uint32, repeatCount uint32, callback func()) *TimerEvent {
	// minimum unit 1ms
	if delay < 1 {
		callback()
		return nil
	}

	// create event object
	event := &TimerEvent{
		callBack:    callback,
		delay:       delay,
		repeatCount: repeatCount,
		closeChan:   make(chan int),
	}

	// start timer
	go startTicker(event)

	// return
	return event
}

func startTicker(event *TimerEvent) {
	defer errutils.PrintPanicStackError()
	event.ticker = time.NewTicker(time.Duration(event.delay) * time.Millisecond)
	for {
		select {
		case <-event.ticker.C:
			event.callBack()
			if event.repeatCount > 0 {
				event.repeatCount -= 1
				if event.repeatCount == 0 {
					Remove(event)
				}
			}
		case <-event.closeChan:
			return
		}
	}
}
