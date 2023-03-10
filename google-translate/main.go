package main

import (
	"flag"
	"fmt"
	"github.com/mireaulf/Learn-GO-by-Building-Three-Simple-Golang-Projects/google-translate/domain"
	"github.com/mireaulf/Learn-GO-by-Building-Three-Simple-Golang-Projects/google-translate/service"
	"strings"
)

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

	svc := service.TranslatorService{}

	p, c := domain.MakeProducerConsumer(requests)

	c.Start()
	p.Start(&svc)
	p.Wait()
	c.Wait()
}
