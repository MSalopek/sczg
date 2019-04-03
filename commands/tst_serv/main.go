package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sczg/web"
	"time"
)

/* Barebones http server serving html templates
* used when running end-to-end tests
* not fit for any other purpose
 */

func serveTestTable(w http.ResponseWriter, r *http.Request) {
	a := make([]web.Advert, 0)
	a = append(a, web.Advert{
		Source:  "TST",
		AdNum:   "TT1",
		Cat:     "NONE",
		Desc:    "SOME DESC",
		Link:    "example@mail.com",
		Contact: "01/12331",
		Rate:    "0.00",
		Date:    "2016-05-06",
	})
	a = append(a, web.Advert{
		Source:  "TST",
		AdNum:   "TT2",
		Cat:     "NONE",
		Desc:    "SOME DESC that is way longer than the previous",
		Link:    "example2@mail.com",
		Contact: "01/1s56a331; 1568132/65",
		Rate:    "30.00",
		Date:    "2016-05-06",
	})
	select {
	case <-r.Context().Done():
		fmt.Println("CLIENT DISCONNECTED FROM INDX PG")
		return

	default:
		t, err := template.ParseFiles("./commands/tst_serv/templates/agregator.html")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("TEST TABLE RAN")
		err = t.Execute(w, a)
		if err != nil {
			log.Fatal("FAILED EXEC TEMPL:", err)
		}
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	// causing client to time out
	select {
	case <-r.Context().Done():
		fmt.Println("CLIENT DISCONNECTED FROM INDX PG")
		return

	case <-time.After(time.Second * 1): // this should cause the client to cancel req
		log.Println(r.Header)
		t, err := template.ParseFiles("./commands/tst_serv/templates/pg_indx_.html")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("INDEX RAN")
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal("FAILED EXEC TEMPL:", err)
		}
	}
}

func serveRaz(w http.ResponseWriter, r *http.Request) {
	// causing client to time out
	select {
	case <-r.Context().Done():
		fmt.Println("CLIENT DISCONNECTED FROM RAZNO PG")
		return
	default:
		t, err := template.ParseFiles("./commands/tst_serv/templates/ponuda_razno_test.html")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("RAZNO RAN")
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal("FAILED EXEC TEMPL:", err)
		}

	}
}

func serveAdmin(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./commands/tst_serv/templates/ponuda_admin_test.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ADMIN RAN")
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("FAILED EXEC TEMPL:", err)
	}
}

func serveSkladiste(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./commands/tst_serv/templates/ponuda_skladiste_test.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("SKLADISTE RAN")
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("FAILED EXEC TEMPL:", err)
	}
}

func serveTurizam(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./commands/tst_serv/templates/ponuda_turizam_test.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("PRODAJA RAN")
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("FAILED EXEC TEMPL:", err)
	}
}

func serveProdaja(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./commands/tst_serv/templates/ponuda_prodaja_test.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("PRODAJA RAN")
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("FAILED EXEC TEMPL:", err)
	}
}

func serveProiz(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./commands/tst_serv/templates/ponuda_proizvodnja_test.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("PROIZ RAN")
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("FAILED EXEC TEMPL:", err)
	}
}

func servePromo(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./commands/tst_serv/templates/ponuda_promo_test.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("PROMO RAN")
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("FAILED EXEC TEMPL:", err)
	}
}

func serveFiz(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./commands/tst_serv/templates/ponuda_fizicki_test.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("FIZ RAN")
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("FAILED EXEC TEMPL:", err)
	}
}

func serveCisc(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./commands/tst_serv/templates/ponuda_ciscenje_test.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("CISC RAN")
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("FAILED EXEC TEMPL:", err)
	}
}

func serveGosti(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./commands/tst_serv/templates/ponuda_ugostitelji_test.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GOSTI RAN")
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("FAILED EXEC TEMPL:", err)
	}
}

func main() {
	http.HandleFunc("/index", serveIndex)
	http.HandleFunc("/razno", serveRaz)
	http.HandleFunc("/admin", serveAdmin)
	http.HandleFunc("/skladiste", serveSkladiste)
	http.HandleFunc("/prodaja", serveProdaja)
	http.HandleFunc("/ciscenje", serveCisc)
	http.HandleFunc("/proizvodnja", serveProiz)
	http.HandleFunc("/promo", servePromo)
	http.HandleFunc("/fizicki", serveFiz)
	http.HandleFunc("/ugostitelji", serveGosti)
	http.HandleFunc("/turizam", serveTurizam)
	http.HandleFunc("/tables", serveTestTable)
	fmt.Println("PORT :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
