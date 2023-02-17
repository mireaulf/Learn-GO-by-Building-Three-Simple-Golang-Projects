package main

import (
	"flag"
	"fmt"
	"github.com/mireaulf/Learn-GO-by-Building-Three-Simple-Golang-Projects/google-translate/domain"
	"github.com/mireaulf/Learn-GO-by-Building-Three-Simple-Golang-Projects/google-translate/service"
	"strings"
	"sync"
)

func produce(requests []domain.Request, wg *sync.WaitGroup, ch *chan *domain.Request) {
	for i := range requests {
		req := requests[i]
		go service.Translate(&req, wg, *ch)
		fmt.Println(fmt.Sprintf("SENT [%v]->[%v] '%v' ", req.SrcLang, req.TgtLang, req.SrcText))
	}
	wg.Done()
}

func consume(wg *sync.WaitGroup, ch *chan *domain.Request) {
	for resp := range *ch {
		fmt.Println(fmt.Sprintf("RECEIVED [%v]->[%v] '%v': '%v'", resp.SrcLang, resp.TgtLang, resp.SrcText, resp.TgtText))
	}
	wg.Done()
}

func main() {
	var srcLang string
	var targetsStr string
	var srcText string

	flag.StringVar(&srcLang, "s", "en", "Source lang [en]")
	flag.StringVar(&targetsStr, "t", "fr", "Target lang [fr]")
	flag.StringVar(&srcText, "st", "", "Text to translate")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	targets := strings.Split(targetsStr, ",")
	requests := make([]domain.Request, 0)
	for _, target := range targets {
		requests = append(requests, domain.Request{
			SrcLang: srcLang,
			TgtLang: target,
			SrcText: srcText,
		})
	}

	ch := make(chan *domain.Request)

	wgConsume := sync.WaitGroup{}
	wgConsume.Add(1)
	go consume(&wgConsume, &ch)

	wgProduce := sync.WaitGroup{}
	wgProduce.Add(1)
	go produce(requests, &wgProduce, &ch)
	wgProduce.Wait()

	close(ch)

	wgConsume.Wait()
}
