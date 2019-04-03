package commands

import (
	"html/template"
	"net/http"
	"sczg/config"
	"sczg/dbutil"

	log "github.com/sirupsen/logrus"
)

/*CURRENT Pages
// Svi
// NOVO
// Istaknuto
// Administracija
// Promocije
// Trgovina
// Ugostiteljstvo
// Turizam
// Čišćenje
// Razno
// Proizvodnja
// Skladišta
// Fizički
*/

func handlerFuncFactory(s *dbutil.Storage, resource string, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			log.Infof("client disconnected from %v", resource)
			return
		default:
			// TODO fix err handling
			ads, err := s.FetchActiveByCategory(key)
			if err != nil {
				log.Error(err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			t, err := template.ParseFiles("./templates/ad_page.html")
			if err != nil {
				log.Error(err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			log.Infof("%v handler executed", key)
			err = t.Execute(w, ads)
			if err != nil {
				log.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	}
}

func serveAll(s *dbutil.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			log.Infof("client disconnected from all")
			return
		default:
			ads, err := s.FetchAllActive()
			if err != nil {
				log.Error(err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			t, err := template.ParseFiles("./templates/ad_page.html")
			if err != nil {
				log.Error(err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			log.Println("serveAllExecuted")
			err = t.Execute(w, ads)
			if err != nil {
				log.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	}
}

func serveNew(s *dbutil.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			log.Infof("client disconnected from new")
			return
		default:
			ads, err := s.FetchFreshActive()
			if err != nil {
				log.Error(err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			t, err := template.ParseFiles("./templates/ad_page.html")
			if err != nil {
				log.Error(err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			log.Println("serveNewExecuted")
			err = t.Execute(w, ads)
			if err != nil {
				log.Error(err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}
	}
}

// StartServer initiates http server
// in a given environment env
func StartServer(env *config.Env) {
	defer env.DB.Close()
	reg := map[string]string{
		"/istaknuto":   "Index",
		"/promo":       "Promo",
		"/trgovina":    "Trgovina",
		"/ugostitelji": "Ugostitelji",
		"/ciscenje":    "Ciscenje",
		"/proizvodnja": "Proizvodnja",
		"/turizam":     "Turizam",
		"/fizicki":     "Fizicki",
		"/razno":       "Razno",
		"/admin":       "Admin",
		"/skladiste":   "Skladiste",
	}
	// register handlers
	for k, v := range reg {
		http.HandleFunc(k, handlerFuncFactory(env.DB, k, v))
	}
	http.HandleFunc("/all", serveAll(env.DB))
	http.HandleFunc("/new", serveNew(env.DB))
	log.Infof("app server running on port: %v", env.Cfg.Port)
	log.Fatal(http.ListenAndServe(env.Cfg.Port, nil))
}
