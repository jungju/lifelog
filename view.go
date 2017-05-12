package main

import (
	"html/template"
	"net/http"

	"github.com/Sirupsen/logrus"
)

const (
	tmplPath = "template/"
)

func handlerHomeView(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		tmplPath + "home.tmpl",
	)
	if err != nil {
		logrus.WithError(err).Error("Failed ParseFiles")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := r.URL.Query().Get("token")
	vaildToken := false
	urlToken := ":token"
	if token != "" {
		_, err = jDB.GetJawbone(token)
		if err == nil {
			vaildToken = true
			urlToken = token
		} else {
			logrus.WithError(err).Error("Failed GetJawbone")
			vaildToken = false
		}
	}

	var data = map[string]interface{}{
		"token":      token,
		"urlToken":   urlToken,
		"vaildToken": vaildToken,
	}
	if err := t.ExecuteTemplate(w, "events", data); err != nil {
		logrus.WithError(err).Error("Failed ExecuteTemplate")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, nil); err != nil {
		logrus.WithError(err).Error("Failed Execute")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handlerEventsView(w http.ResponseWriter, r *http.Request) {
}

func handlerTokensView(w http.ResponseWriter, r *http.Request) {
}
