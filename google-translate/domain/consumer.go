package domain

import (
	"fmt"
	"sync"
)

type Consumer struct {
	ch *chan *Request
	wg sync.WaitGroup
}

func (c *Consumer) Start() {
	c.wg.Add(1)
	go func() {
		for resp := range *c.ch {
			fmt.Println(fmt.Sprintf("RECEIVED [%v]->[%v] '%v': '%v'", resp.SrcLang, resp.TgtLang, resp.SrcText, resp.TgtText))
		}
		c.wg.Done()
	}()
}

func (c *Consumer) Wait() {
	c.wg.Wait()
}
