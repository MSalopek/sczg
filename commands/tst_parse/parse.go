package main

import (
	"fmt"
	"log"
	"os"

	"github.com/antchfx/htmlquery"

	"sczg/web"
)

func main() {
	r := web.PrepareRgx()
	f, err := os.Open("/home/ms/go/src/sczg/commands/tst_serv/templates/ponuda_ugostitelji_test.html")
	if err != nil {
		log.Fatal(err)
	}
	p, err := htmlquery.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	ads, err := web.ParseDedicatedPg(p, "TEST", r)
	fmt.Println(len(ads))
	for _, val := range ads {
		fmt.Println(val.AdNum)
		// fmt.Println(val.AdNum)
		// fmt.Printf("%v\n", val.String())
	}
}
