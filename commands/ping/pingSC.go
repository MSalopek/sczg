package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	c := &http.Client{
		Timeout: time.Second * 10,
	}
	rq, err := http.NewRequest("GET", "http://www.sczg.unizg.hr/student-servis/vijest/2015-04-14-razni-poslovi/", nil)
	rq.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:62.0) Gecko/20100101 Firefox/62.0")
	rq.Header.Set("Content-Type", "text/html; charset=utf-8")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := c.Do(rq)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RESP H", resp.Header, "RESP S", resp.Status)
	fmt.Println(string(b))
	fmt.Printf("REQ HEAD\n %+v, REQUEST: \n%+v\nAGENT: %v\n", rq.Header, rq, rq.UserAgent())
}
