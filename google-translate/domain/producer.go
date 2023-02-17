package domain

import (
	"fmt"
	"sync"
)

type Producer struct {
	requests []Request
	ch       *chan *Request
	wg       sync.WaitGroup
}

func (p *Producer) Start(tr Translator) {
	p.wg.Add(1)
	go func() {
		for i := range p.requests {
			req := p.requests[i]
			go tr.Translate(&req, &p.wg, *p.ch)
			fmt.Println(fmt.Sprintf("SENT [%v]->[%v] '%v' ", req.SrcLang, req.TgtLang, req.SrcText))
		}
		p.wg.Done()
	}()
}

func (p *Producer) Wait() {
	p.wg.Wait()
	close(*p.ch)
}
