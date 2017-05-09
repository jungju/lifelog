package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"

	"golang.org/x/oauth2"

	goji "goji.io"

	"github.com/Sirupsen/logrus"
	"goji.io/pat"
)

var (
	dbName string
	port   string
	host   string
	jDB    *jawboneDB
)

var (
	conf = &oauth2.Config{
		RedirectURL:  "http://life.jjgo.kr/up/callback",
		ClientID:     os.Getenv("JAWBONE_KEY"),
		ClientSecret: os.Getenv("JAWBONE_SECRET"),
		Scopes:       []string{"basic_read", "extended_read", "location_read", "friends_read", "mood_read", "mood_write", "move_read", "move_write", "sleep_read", "sleep_write", "meal_read", "meal_write", "weight_read", "weight_write", "generic_event_read", "generic_event_write", "heartrate_read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://jawbone.com/auth/oauth2/auth",
			TokenURL: "https://jawbone.com/auth/oauth2/token",
		},
	}
)

func init() {
	dbName = os.Getenv("DB_NAME")
	port = os.Getenv("WEB_PORT")
	host = os.Getenv("HOST")
	if port == "" {
		port = "8080"
	}
	if host == "" {
		host = "localhost"
	}
	if dbName == "" {
		dbName = "db/life.db"
	}
}

func main() {
	logrus.Info("starting life-server..")
	logrus.SetLevel(logrus.DebugLevel)

	var err error
	jDB = &jawboneDB{}
	if jDB.DB, err = bolt.Open(dbName, 0600, &bolt.Options{Timeout: 2 * time.Second}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/up/redirect"), handlerRedirectToLogin)
	mux.HandleFunc(pat.Get("/up/callback"), handlerCallback)
	mux.HandleFunc(pat.Get("/up/event/:token/water"), handlerWater)
	mux.HandleFunc(pat.Get("/up/event/:token/coffee"), handlerCoffee)
	mux.HandleFunc(pat.Get("/up/event/:token/pooh"), handlerPooh)
	mux.HandleFunc(pat.Get("/up/event/:token/urine"), handlerUrine)
	mux.HandleFunc(pat.Get("/up/event/:token/custom"), handlerCustom)
	mux.HandleFunc(pat.Get("/up/tokens"), handlerTokens)
	mux.HandleFunc(pat.Get("/"), handlerHomeView)
	//mux.HandleFunc(pat.Get("/tokens"), handlerTokensView)
	//mux.HandleFunc(pat.Get("/:token/events"), handlerEventsView)
	mux.Handle(pat.Get("/statics/*"), http.StripPrefix("/statics/", http.FileServer(http.Dir("statics"))))
	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
