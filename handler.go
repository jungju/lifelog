package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"

	"time"

	"strings"

	"github.com/Sirupsen/logrus"
	"goji.io/pat"
)

func handlerRedirectToLogin(w http.ResponseWriter, r *http.Request) {
	AuthCodeOption := oauth2.SetAuthURLParam("access_type", "code")
	url := conf.AuthCodeURL("state", AuthCodeOption)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handlerCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	queryString := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=authorization_code&code=%s", conf.ClientID, conf.ClientSecret, code)
	req := httpRequestParams{
		URL:    fmt.Sprintf("%s?%s", conf.Endpoint.TokenURL, queryString),
		Method: "POST",
	}

	res, err := req.request()
	if err != nil {
		logrus.WithError(err).Error("Failed callback")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resToken := &token{}
	if err = json.Unmarshal(res.Body, resToken); err != nil {
		logrus.WithError(err).Error("Failed Unmarshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	createJawbone := &jawbone{Token: resToken}
	if err := createJawbone.save(); err != nil {
		logrus.WithError(err).Error("Failed save")
		w.WriteHeader(http.StatusInternalServerError)
	}
	http.Redirect(w, r, fmt.Sprintf("/?token=%s", createJawbone.ID), http.StatusTemporaryRedirect)
}

func makeJawbone(r *http.Request) (*jawbone, error) {
	id := pat.Param(r, "token")
	return jDB.GetJawbone(id)
}

func handlerTokens(w http.ResponseWriter, r *http.Request) {
	jawbones, err := jDB.ListJawbons()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	for _, jawbone := range jawbones {
		fmt.Fprintln(w, jawbone.ID, jawbone.Token)
	}
	w.WriteHeader(http.StatusOK)
}

func handlerWater(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cups := convToIntFromQuery(r, "cups", -1)
	if cups == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Require cups")
		return
	}

	if err := jawboneClient.eventWater(time.Now(), int(cups)); err != nil {
		logrus.WithError(err).Error("Failed drinkWater")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handlerIcedAmericano(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := jawboneClient.eventIcedAmericano(time.Now()); err != nil {
		logrus.WithError(err).Error("Failed eventIcedAmericano")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handlerIcedLatte(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := jawboneClient.eventIcedLatte(time.Now()); err != nil {
		logrus.WithError(err).Error("Failed eventIcedLatte")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handlerPooh(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	poohType := convToIntFromQuery(r, "poohType", -1)
	if poohType == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Require poohType")
		return
	}
	pain := convToBoolFromQuery(r, "pain", false)
	constipation := convToBoolFromQuery(r, "constipation", false)
	blood := convToBoolFromQuery(r, "blood", false)

	if err := jawboneClient.eventPooh(time.Now(), poohType, pain, constipation, blood); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handlerUrine(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	peeType := convToIntFromQuery(r, "peeType", -1)
	if peeType == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Require peeType")
		return
	}
	color := convToIntFromQuery(r, "color", 1)
	blood := convToBoolFromQuery(r, "blood", false)

	if err := jawboneClient.eventUrine(time.Now(), peeType, color, blood); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func handlerMigraine(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	direction := convToIntFromQuery(r, "direction", -1)
	if direction == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Require direction")
		return
	}

	if err := jawboneClient.eventMigraine(time.Now(), direction); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func handlerIndigestion(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	organ := convToIntFromQuery(r, "organ", -1)
	if organ == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Require peeType")
	}

	if err := jawboneClient.eventIndigestion(time.Now(), organ); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func equalMealItemType(itemType string, name string) bool {
	if strings.Index(name, fmt.Sprintf("%s_", itemType)) == 0 {
		return true
	}
	return false
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
