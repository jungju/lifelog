package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/Sirupsen/logrus"
)

const (
	tmplPath = "template/"
)

func handlerHomeView(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		tmplPath + "home.html",
	)
	if err != nil {
		logrus.WithError(err).Error("Failed ParseFiles")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
		return
	}
	token := r.URL.Query().Get("token")

	vaildToken := true
	_, err = jDB.GetJawbone(token)
	if err != nil {
		token = ""
		vaildToken = false
	}

	urlToken := ""
	if token != "" {
		urlToken = token
	} else {
		urlToken = ":token"
	}

	var data = map[string]interface{}{
		"token":      token,
		"urlToken":   urlToken,
		"vaildToken": vaildToken,
	}
	if err := t.ExecuteTemplate(w, "events", data); err != nil {
		logrus.WithError(err).Error("Failed ExecuteTemplate")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
		return
	}
	if err := t.Execute(w, nil); err != nil {
		logrus.WithError(err).Error("Failed Execute")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
		return
	}
}

func handlerEventsView(w http.ResponseWriter, r *http.Request) {
}

func handlerTokensView(w http.ResponseWriter, r *http.Request) {
}
