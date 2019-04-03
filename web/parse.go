package web

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

// Advert stores ad data
type Advert struct {
	Date    string
	Source  string
	AdNum   string
	Cat     string
	Desc    string
	Link    string
	Contact string
	Rate    string
}

func (a *Advert) String() string {
	return fmt.Sprintf("NUM: %v\nCAT: %v\nSRC: %v\nDAT: %v\nRTE: %v\nLNK: %v\nCON: %v\nDSC: %v",
		a.AdNum, a.Cat, a.Source, a.Date, a.Rate, a.Link, a.Contact, a.Desc)
}

type regExes struct {
	Emails    *regexp.Regexp
	Rate      *regexp.Regexp
	RateFallB *regexp.Regexp
	Phone     *regexp.Regexp
}

// PrepareRgx compiles regExes
func PrepareRgx() *regExes {
	return &regExes{
		Emails:    regexp.MustCompile("[a-zA-Z0-9-_.]+@[a-zA-Z0-9-_.]+"),
		Rate:      regexp.MustCompile(`(\d{2,3}[\,\.]{0,1}\d{0,2})(kn| kuna|kuna| kn| KN| Kn|Kn| kune){1}`),
		RateFallB: regexp.MustCompile(`(\d{2}\ )(kn|kuna|Kn|KN])`),
		Phone:     regexp.MustCompile(`(\d{2,3}[\/|\ |\-|\\]{0,2}\d+[\-|\ ]{0,1}\d{3,6})`),
	}
}

// ParseFPjobBoxes parses short descriptions of
// front page adverts. Each front page advert has
// a full page dedicated to it specified in the ad.Link.
// Full ad info is gathered by other means.
// see:  web.ParseFullPageAd and fetcher service.
func ParseFPjobBoxes(doc *html.Node) ([]Advert, error) {
	adverts := make([]Advert, 0)
	for _, n := range htmlquery.Find(doc, `//div[@id="mainContent"]/div[@class="newsItem"]/div[@class="jobBox"]`) {
		txt := htmlquery.InnerText(n)
		if len(txt) <= 1 {
			continue
		}
		trimTxt := strings.Replace(txt, "\n", "", 10)
		txtElems := strings.Split(trimTxt, "    ") // 4 blanklines
		cleanElems := make([]string, 0)
		for _, elem := range txtElems {
			if len(elem) > 1 && !strings.Contains(elem, "\n") {
				cleanElems = append(cleanElems, strings.Trim(elem, " "))
			}
		}
		linkNode := htmlquery.FindOne(n, `//a[@class="bLink"]/@href`)
		link := htmlquery.SelectAttr(linkNode, "href")
		if len(cleanElems) > 3 {
			return []Advert{}, fmt.Errorf("element list longer than expected {%#v} len shoud be 3", cleanElems)
		}
		adNum := getjBoxAdNumHelper(link)
		ad := Advert{
			Date:    time.Now().Local().Format("2006-01-02 15:04:05"),
			Source:  "jBox",
			AdNum:   "I" + adNum,
			Cat:     "Index",
			Desc:    cleanElems[1] + ": " + cleanElems[2],
			Link:    "http://www.sczg.unizg.hr" + link,
			Contact: "--",
			Rate:    "00.00",
		}
		adverts = append(adverts, ad)
	}
	return adverts, nil
}

// ParseFullPageAd parses a single full page ad.
// Short descriptions of full page ads are handled in ParseFPjBoxes.
func ParseFullPageAd(doc *html.Node, adnum string, reg *regExes) (Advert, error) {
	var ad Advert
	html := htmlquery.FindOne(doc, `//div[@id="mainContent"]/div[@class="newsItem"]/div[@class="content"]`)
	txt := []byte(strings.Join(strings.Fields(htmlquery.InnerText(html)), " "))
	if len(txt) == 0 {
		return Advert{}, fmt.Errorf("no text found while parsing full page ad {%v}", adnum)
	}
	ad.Date = time.Now().Local().Format("2006-01-02 15:04:05")
	ad.Source = "jBox"
	ad.Cat = "Index"
	ad.AdNum = adnum
	ad.Desc = string(txt)
	advertDetailsHelper(txt, reg, &ad)

	log.Infof("finished parsing full page ad {%v}", adnum)
	return ad, nil
}

// ParseDedicatedPg parses a single page
// containing multiple ads of a specific category.
func ParseDedicatedPg(doc *html.Node, cat string, reg *regExes) ([]Advert, error) {
	adverts := make([]Advert, 0)
	for i, n := range htmlquery.Find(doc, `//div[@id="mainContent"]/div[@class="newsItem"]/div[@class="content"]//p`) {
		tt := htmlquery.InnerText(n)
		if len(tt) <= 1 || i == 0 {
			continue
		}
		// EXAMPLE: 5730/ Konobarski poslovi u caffe baru u sklopu Supermarketa u Rovinju (...)
		// Number before "/" not unique, adding Cat[0] to achieve Adnum uniqueness for db reasons
		numNtext := strings.SplitN(tt, "/", 2)
		num := string(cat[0]) + numNtext[0]
		ad := Advert{
			Date:   time.Now().Local().Format("2006-01-02 15:04:05"), // go date formatting gotchas...
			Source: "Dedicated",
			Cat:    cat,
			AdNum:  num,
			Desc:   numNtext[1],
		}
		advertDetailsHelper([]byte(numNtext[1]), reg, &ad)
		adverts = append(adverts, ad)
	}
	log.Infof("found ads in %v:  %v", cat, len(adverts))
	return adverts, nil
}

func advertDetailsHelper(b []byte, reg *regExes, ad *Advert) {
	emailsRaw := reg.Emails.FindAll(b, 2)
	emails := ""
	for j, k := range emailsRaw {
		if j == 0 {
			emails = string(k)
		} else {
			emails += "; " + string(k)
		}
	}
	if emails == "" {
		emails = "--"
	}
	phoneRaw := reg.Phone.FindAll(b, 2)
	phone := "--"
	for j, k := range phoneRaw {
		if j == 0 {
			phone = string(k)
		} else {
			phone += "; " + string(k)
		}
	}
	rate := string(reg.Rate.Find(b))
	if rate == "" {
		rate = string(reg.RateFallB.Find(b))
		// if there really is no such thing
		if rate == "" {
			rate = "00.00"
		}
	}
	rate = formatRatesHelper(rate)

	ad.Link = emails
	ad.Contact = phone
	ad.Rate = rate
}

func getjBoxAdNumHelper(s string) string {
	// EXPAMPLE STRING: /student-servis/poslovi/11384/
	numStr := s[24:]
	numStr = strings.Trim(numStr, "/")
	return numStr
}

// dealing with human inputted strings
// is always FUN!
func formatRatesHelper(s string) string {
	var rate string
	switch {
	case strings.HasSuffix(s, ".00 kn"):
		rate = strings.TrimRight(s, " kn")
	case strings.HasSuffix(s, ".00 kuna"):
		rate = strings.TrimRight(s, " kuna")
	case strings.HasSuffix(s, ".00 kune"):
		rate = strings.TrimRight(s, " kune")
	case strings.HasSuffix(s, ",00 kune"):
		intPart := s[:strings.Index(s, ",")]
		rate = intPart + ".00"
	case strings.HasSuffix(s, ",00 kn"):
		intPart := s[:strings.Index(s, ",")]
		rate = intPart + ".00"
	case strings.HasSuffix(s, ",00 kuna"):
		intPart := s[:strings.Index(s, ",")]
		rate = intPart + ".00"
	case strings.HasSuffix(s, " k"):
		intPart := s[:strings.Index(s, "k")]
		rate = intPart + ".00"
	case strings.HasSuffix(s, " kn"):
		intPart := s[:strings.Index(s, " ")]
		rate = intPart + ".00"
	case strings.HasSuffix(s, " Kn"):
		intPart := s[:strings.Index(s, " ")]
		rate = intPart + ".00"
	case strings.HasSuffix(s, " KN"):
		intPart := s[:strings.Index(s, " ")]
		rate = intPart + ".00"
	case strings.HasSuffix(s, ",00kn"):
		intPart := s[:strings.Index(s, ",")]
		rate = intPart + ".00"
	case strings.HasSuffix(s, "kn"):
		intPart := s[:strings.Index(s, "k")]
		rate = intPart + ".00"
	case strings.HasSuffix(s, " "):
		intPart := s[:strings.Index(s, " ")]
		rate = intPart + ".00"
	case s == "00.00":
		rate = s
	}
	return rate
}
