package web

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

type Result struct {
	Ads []Advert
	Err error
}

func fetchURL(cli *http.Client, url string) (*html.Node, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:62.0) Gecko/20100101 Firefox/62.0")
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	return html.Parse(r)
}

// ProcessPage fetches URL and parses data
// parsed data is forwarded into a chan of
// type Result for further manipulation
func ProcessPage(cli *http.Client, url string, ch chan<- Result, key string, reg *regExes) {
	doc, err := fetchURL(cli, url)
	if err != nil {
		Err := fmt.Errorf("Error processing page {%v}\n DETAILS: {%v}", key, err)
		ch <- Result{[]Advert{}, Err}
		return
	}

	if key == "Index" {
		ads, err := ParseFPjobBoxes(doc)
		if err != nil {
			Err := fmt.Errorf("Error processing page {%v}\n DETAILS: {%v}", key, err)
			ch <- Result{[]Advert{}, Err}
		}
		ch <- Result{ads, nil}
	} else {
		ads, err := ParseDedicatedPg(doc, key, reg)
		if err != nil {
			Err := fmt.Errorf("Error processing page {%v}\n DETAILS: {%v}", key, err)
			ch <- Result{[]Advert{}, Err}
		}
		ch <- Result{ads, nil}
	}
}
