package main

import (
	"encoding/xml"
	"fmt"
	"github.com/mireaulf/Learn-GO-by-Building-Three-Simple-Golang-Projects/google-trends/domain"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	region := "US"
	if len(os.Args) > 1 {
		region = os.Args[1]
	}
	fmt.Printf("fetching google trends (region:%s)\n", region)
	rss := readGoogleTrends(region)

	fmt.Printf("TRENDS:\n%s", rss)
}

func readGoogleTrends(region string) domain.RSS {
	var rss domain.RSS
	data := getGoogleTrendsRss(region)
	if err := xml.Unmarshal(data, &rss); err != nil {
		fmt.Printf("Error unmarshaling xml: %v\n", err)
	}
	return rss
}

func getGoogleTrendsRss(region string) []byte {
	url := fmt.Sprintf("https://trends.google.com/trends/trendingsearches/daily/rss?geo=%s", region)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error getting rss: %v", err)
		return nil
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		return nil
	}
	return content
}
