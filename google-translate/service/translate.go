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

type TranslatorService struct{}

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
		return fmt.Sprintf("Error unmarshalling level 0: %T: '%v'", resp, resp)
	}
	for i := 0; i < 2; i++ {
		if data, ok = data[0].([]any); !ok {
			return fmt.Sprintf("Error unmarshalling level %v: %T: '%v'", i, data[0], data[0])
		}
	}
	s, ok := data[0].(string)
	if !ok {
		return fmt.Sprintf("Error unmarshalling string: %T : '%v'", data[0], data[0])
	}
	return s
}

func (*TranslatorService) Translate(req *domain.Request, wg *sync.WaitGroup, ch chan *domain.Request) {
	wg.Add(1)
	defer wg.Done()
	httpReq, err := buildQuery(req)
	if err != nil {
		req.TgtText = fmt.Sprintf("Error building new request: %v", err)
		ch <- req
		return
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		req.TgtText = fmt.Sprintf("Error making request: %v", err)
		ch <- req
		return
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		req.TgtText = fmt.Sprintf("Request failed with code %v", httpResp.StatusCode)
		ch <- req
		return
	}
	b, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		req.TgtText = fmt.Sprintf("Error reading body %v", err)
		ch <- req
		return
	}
	req.TgtText = unmarshalResponse(b)
	ch <- req
}
