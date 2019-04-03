package main

import (
	"fmt"
	"sczg/config"
)

func main() {
	c, _ := config.InitCfg("./config/config.yaml")
	m := c.MapURLs()
	for k, v := range m {
		fmt.Println(k, ":", v)
	}
}
