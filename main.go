package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"

	"golang.org/x/oauth2"

	goji "goji.io"

	"strconv"

	"log"

	"goji.io/pat"
)

const version = "v0.0.0"

var (
	db     *bolt.DB
	dbName string
	port   string
)

var (
	conf = &oauth2.Config{
		RedirectURL:  "https://jjgo.kr/apps/up/redirect",
		ClientID:     os.Getenv("JAWBONE_KEY"),
		ClientSecret: os.Getenv("JAWBONE_SECRET"),
		Scopes:       []string{"basic_read", "extended_read", "location_read", "friends_read", "mood_read", "mood_write", "move_read", "move_write", "sleep_read", "sleep_write", "meal_read", "meal_write", "weight_read", "weight_write", "generic_event_read", "generic_event_write", "heartrate_read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://jawbone.com/auth/oauth2/auth",
			TokenURL: "https://jawbone.com/auth/oauth2/token",
		},
	}
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	AuthCodeOption := oauth2.SetAuthURLParam("access_type", "code")
	url := conf.AuthCodeURL("state", AuthCodeOption)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func callback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	exch, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
	}

	j := jawbone{AccessToken: exch.AccessToken}
	if err := j.save(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
	}
	fmt.Fprint(w, j.ID)
}

func makeJawbone(r *http.Request) (*jawbone, error) {
	id := pat.Param(r, "id")
	return getJawbone(id)
}

func handlerWater(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
	}

	cups := convToIntFromQuery(r, "cups", -1)
	if cups == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Require cups")
	}
	if cups == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
	}

	if err := jawboneClient.drinkWater(int(cups)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
	}
}

func handlerDump(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
	}

	dumptype := convToIntFromQuery(r, "type", -1)
	if dumptype == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Require dumptype")
	}
	pain := convToBoolFromQuery(r, "pain", false)
	constipation := convToBoolFromQuery(r, "constipation", false)
	blood := convToBoolFromQuery(r, "blood", false)

	if err := jawboneClient.takeDump(dumptype, pain, constipation, blood); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
	}
}

func handlerPee(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
	}

	peeType := convToIntFromQuery(r, "peeType", -1)
	if peeType == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Require peeType")
	}
	blood := convToBoolFromQuery(r, "blood", false)

	if err := jawboneClient.pee(peeType, blood); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
	}
}

func handlerCustom(w http.ResponseWriter, r *http.Request) {
	// jawboneClient, err := makeJawbone(r)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprint(w, "Internal Server Error 500")
	// }

	// if err := jawboneClient.//(); err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprint(w, "Internal Server Error 500")
	// }
}

func convToIntFromQuery(r *http.Request, query string, defaultValue int) int {
	strTarget, err := strconv.ParseInt(r.URL.Query().Get(query), 10, 64)
	if err != nil {
		return defaultValue
	}
	return int(strTarget)
}

func convToBoolFromQuery(r *http.Request, query string, defaultValue bool) bool {
	strTarget := r.URL.Query().Get(query)
	if strTarget == "1" {
		return true
	} else if strTarget == "0" {
		return false
	}
	return defaultValue
}

func init() {
	dbName = os.Getenv("DB_NAME")
	port = os.Getenv("WEB_PORT")
	if port == "" {
		port = "8080"
	}
}

func main() {
	log.Println("starting boltdb-browser..")

	var err error
	if db, err = bolt.Open(dbName, 0600, &bolt.Options{Timeout: 2 * time.Second}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/"), hello)
	mux.HandleFunc(pat.Get("/up/redirect"), redirectToLogin)
	mux.HandleFunc(pat.Post("/up/callback"), callback)
	mux.HandleFunc(pat.Post("/up/action/:id/water"), handlerWater)
	mux.HandleFunc(pat.Post("/up/action/:id/dump"), handlerDump)
	mux.HandleFunc(pat.Post("/up/action/:id/pee"), handlerPee)
	mux.HandleFunc(pat.Post("/up/action/:id/custom"), handlerCustom)

	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
