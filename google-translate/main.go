package main

import (
	"flag"
	"fmt"
	"github.com/mireaulf/Learn-GO-by-Building-Three-Simple-Golang-Projects/google-translate/domain"
	"github.com/mireaulf/Learn-GO-by-Building-Three-Simple-Golang-Projects/google-translate/service"
	"sync"
)

func main() {
	req := domain.Request{}
	flag.StringVar(&req.SrcLang, "s", "en", "Source lang [en]")
	flag.StringVar(&req.TgtLang, "t", "fr", "Target lang [fr]")
	flag.StringVar(&req.SrcText, "st", "", "Text to translate")
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	ch := make(chan string)
	wg := sync.WaitGroup{}
	go service.Translate(&req, &wg, ch)
	select {
	case text := <-ch:
		fmt.Println(fmt.Sprintf("[%v]->[%v] '%v' : %v", req.SrcLang, req.TgtLang, req.SrcText, text))
	}
	wg.Wait()
	close(ch)
}
