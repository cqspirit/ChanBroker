package ChanBroker

import (
	"container/list"
	"errors"
	"time"
)

type Content interface{}

type Subscriber chan Content

type ChanBroker struct {
	regSub      chan Subscriber
	unRegSub    chan Subscriber
	contents    chan Content
	stop        chan bool
	subscribers map[Subscriber]*list.List
	timeout     time.Duration
	cachenum    uint
	timerChan   <-chan time.Time
}

var ErrBrokerExit error = errors.New("ChanBroker exit")
var ErrPublishTimeOut error = errors.New("ChanBroker Pulish Time out")
var ErrRegTimeOut error = errors.New("ChanBroker Reg Time out")
var ErrStopPublishTimeOut error = errors.New("ChanBroker Stop Publish Time out")

func NewChanBroker(timeout time.Duration) *ChanBroker {
	Broker := new(ChanBroker)
	Broker.regSub = make(chan Subscriber)
	Broker.unRegSub = make(chan Subscriber)
	Broker.contents = make(chan Content)
	Broker.stop = make(chan bool)

	Broker.subscribers = make(map[Subscriber]*list.List)
	Broker.timeout = timeout
	Broker.cachenum = 0
	Broker.timerChan = nil
	Broker.run()

	return Broker
}

func (self *ChanBroker) onContentPush(content Content) {
	for sub, clist := range self.subscribers {

		for next := clist.Front(); next != nil; {
			cur := next
			next = cur.Next()
			select {
			case sub <- cur.Value:
				if self.cachenum > 0 {
					self.cachenum--
				}
				clist.Remove(cur)
			default:
				break // block
			}
		}

		len := clist.Len()
		if len == 0 {
			select {
			case sub <- content:
			default:
				clist.PushBack(content)
				self.cachenum++
			}
		} else {
			clist.PushBack(content)
			self.cachenum++
		}
	}

	if self.cachenum > 0 && self.timerChan == nil {
		timer := time.NewTimer(self.timeout)
		self.timerChan = timer.C
	}

}

func (self *ChanBroker) onTimerPush() {
	for sub, clist := range self.subscribers {

		for next := clist.Front(); next != nil; {
			cur := next
			next = cur.Next()
			select {
			case sub <- cur.Value:
				if self.cachenum > 0 {
					self.cachenum--
				}
				clist.Remove(cur)
			default:
				break // block
			}
		}
	}

	if self.cachenum > 0 {
		timer := time.NewTimer(self.timeout)
		self.timerChan = timer.C
	} else {
		self.timerChan = nil
	}
}

func (self *ChanBroker) run() {

	go func() { // Broker Goroutine
		for {
			select {
			case content := <-self.contents:
				self.onContentPush(content)

			case <-self.timerChan:
				self.onTimerPush()

			case sub := <-self.regSub:
				clist := list.New()
				self.subscribers[sub] = clist

			case sub := <-self.unRegSub:
				_, ok := self.subscribers[sub]
				if ok {
					delete(self.subscribers, sub)
					close(sub)
				}

			case <-self.stop:
				for sub := range self.subscribers {
					delete(self.subscribers, sub)
					close(sub)
				}

				return // exit goroutine
			}
		}
	}()
}

func (self *ChanBroker) RegSubscriber(size uint) (Subscriber, error) {
	sub := make(Subscriber, size)

	select {

	case <-time.After(self.timeout):
		return nil, ErrRegTimeOut

	case self.regSub <- sub:
		return sub, nil
	}

}

func (self *ChanBroker) UnRegSubscriber(sub Subscriber) {
	select {
	case <-time.After(self.timeout):
		return

	case self.unRegSub <- sub:
		return
	}

}

func (self *ChanBroker) StopPublish() error {
	select {
	case self.stop <- true:
		return nil

	case <-time.After(self.timeout):
		return ErrStopPublishTimeOut
	}
}

func (self *ChanBroker) PubContent(c Content) error {
	select {
	case <-time.After(self.timeout):
		return ErrPublishTimeOut

	case self.contents <- c:
		return nil
	}

}
