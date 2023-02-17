package domain

import (
	"sync"
)

type Request struct {
	SrcLang string
	TgtLang string
	SrcText string
	TgtText string
}

type Translator interface {
	Translate(req *Request, wg *sync.WaitGroup, ch chan *Request)
}

func MakeProducerConsumer(requests []Request) (*Producer, *Consumer) {
	ch := make(chan *Request)
	return &Producer{
			requests: requests,
			ch:       &ch,
			wg:       sync.WaitGroup{},
		}, &Consumer{
			ch: &ch,
			wg: sync.WaitGroup{},
		}
}
