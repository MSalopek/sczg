package commands

import (
	"net/http"

	"sczg/config"
	"sczg/web"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// TODO DEFINE CRAWL FUNCTION

// StartFetcher starts processing ad pages
// and storing results into database
func StartFetcher(env *config.Env) {
	defer env.DB.Close()
	baseURL := env.Cfg.Base
	var wg sync.WaitGroup
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	regexes := web.PrepareRgx()
	urlMap := env.Cfg.MapURLs()
	results := make(chan web.Result)
	tick := time.NewTicker(time.Second * time.Duration(env.Cfg.Interval)).C

	oneRun := func() {
		wg.Add(len(urlMap))
		for k, v := range urlMap {
			url := baseURL + v
			log.Infof("Processing: %v", k)
			go web.ProcessPage(client, url, results, k, regexes)
		}
		go func() {
			for val := range results {
				if val.Err != nil {
					log.Error(val.Err)
					wg.Done()
				} else {
					env.DB.InsertNewAds(val.Ads)
					wg.Done()
				}

			}
		}()
		wg.Wait()
		log.Infof("finished fetcher")
	}
	oneRun()

	// start subsequent runs on ticker expiration
	// this is blocking
	for {
		select {
		case <-tick:
			log.Infof("running fetcher after interval {%v}", env.Cfg.Interval)
			oneRun()
		}
	}
}
