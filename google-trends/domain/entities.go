package domain

import (
	"encoding/xml"
	"fmt"
)

func leftpad(i int, s string) string {
	out := ""
	for j := 0; j <= i; j++ {
		out += " "
	}
	return fmt.Sprintf("%s%s", out, s)
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel *Channel `xml:"channel"`
}

func (r RSS) String() string {
	return fmt.Sprintf("%s", r.Channel)
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

func (c *Channel) String() string {
	itemsDisplay := ""
	for i, item := range c.Items {
		itemsDisplay += fmt.Sprintf("#%v\n%s\n", i+1, item)
	}
	return fmt.Sprintf("Channel: %s\n%s", c.Title, itemsDisplay)
}

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Traffic string `xml:"approx_traffic"`
	News    []News `xml:"news_item"`
}

func (i Item) String() string {
	newsDisplay := "News:\n"
	for _, news := range i.News {
		newsDisplay += fmt.Sprintf("%s\n", news)
	}
	return fmt.Sprintf("- %s - Traffic: %v (%s)\n%s", i.Title, i.Traffic, i.Link, leftpad(1, newsDisplay))
}

type News struct {
	Headline     string `xml:"news_item_title"`
	HeadlineLink string `xml:"news_item_url"`
}

func (n News) String() string {
	return leftpad(2, fmt.Sprintf("%s: %s", n.Headline, n.HeadlineLink))
}
