package service

import (
	"encoding/json"
	"fmt"
	"github.com/mireaulf/Learn-GO-by-Building-Three-Simple-Golang-Projects/google-translate/domain"
	"io/ioutil"
	"net/http"
	"sync"
)

const TranslateUrl = "https://translate.googleapis.com/translate_a/single"

func buildQuery(req *domain.Request) (*http.Request, error) {
	httpReq, err := http.NewRequest("GET", TranslateUrl, nil)
	if err != nil {
		return nil, err
	}
	httpQuery := httpReq.URL.Query()
	httpQuery.Add("client", "gtx")
	httpQuery.Add("sl", req.SrcLang)
	httpQuery.Add("tl", req.TgtLang)
	httpQuery.Add("dt", "t")
	httpQuery.Add("q", req.SrcText)
	httpReq.URL.RawQuery = httpQuery.Encode()
	return httpReq, nil
}

func unmarshalResponse(b []byte) string {
	var resp interface{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return fmt.Sprintf("Error unmarshaling response %v '%v'", err, string(b))
	}
	data, ok := resp.([]any)
	if !ok {
		return fmt.Sprintf("Error type assert level 0: %T: '%v'", resp, resp)
	}
	for i := 0; i < 2; i++ {
		if data, ok = data[0].([]any); !ok {
			return fmt.Sprintf("Error type assert level %v: %T: '%v'", i, data[0], data[0])
		}
	}
	s, ok := data[0].(string)
	if !ok {
		return fmt.Sprintf("Error unmarshalling string: %T : '%v'", data[0], data[0])
	}
	return s
}

func Translate(req *domain.Request, wg *sync.WaitGroup, ch chan string) {
	wg.Add(1)
	defer wg.Done()

	httpReq, err := buildQuery(req)
	if err != nil {
		ch <- fmt.Sprintf("Error building new request: %v", err)
		return
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		ch <- fmt.Sprintf("Error making request: %v", err)
		return
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		ch <- fmt.Sprintf("Request failed with code %v", httpResp.StatusCode)
		return
	}
	b, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		ch <- fmt.Sprintf("Error reading body %v", err)
		return
	}
	ch <- unmarshalResponse(b)
}
