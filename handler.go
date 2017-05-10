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
		fmt.Fprint(w, "Internal Server Error 500")
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
		fmt.Fprint(w, "Internal Server Error 500")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handlerCoffee(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
		return
	}

	cups := convToIntFromQuery(r, "cups", -1)
	if cups == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Require cups")
		return
	}

	if err := jawboneClient.eventCoffee(time.Now()); err != nil {
		logrus.WithError(err).Error("Failed drinkWater")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handlerPooh(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
		return
	}

	dumptype := convToIntFromQuery(r, "type", -1)
	if dumptype == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Require dumptype")
	}
	pain := convToBoolFromQuery(r, "pain", false)
	constipation := convToBoolFromQuery(r, "constipation", false)
	blood := convToBoolFromQuery(r, "blood", false)

	if err := jawboneClient.eventPooh(time.Now(), dumptype, pain, constipation, blood); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handlerUrine(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
		return
	}

	peeType := convToIntFromQuery(r, "peeType", -1)
	if peeType == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Require peeType")
	}
	blood := convToBoolFromQuery(r, "blood", false)

	if err := jawboneClient.eventUrine(time.Now(), peeType, blood); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
	}

	w.WriteHeader(http.StatusCreated)
}

func handlerMigraine(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
		return
	}

	direction := convToIntFromQuery(r, "direction", -1)
	if direction == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Require direction")
	}

	if err := jawboneClient.eventMigraine(time.Now(), direction); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
	}

	w.WriteHeader(http.StatusCreated)
}

func handlerIndigestion(w http.ResponseWriter, r *http.Request) {
	jawboneClient, err := makeJawbone(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
		return
	}

	organ := convToIntFromQuery(r, "organ", -1)
	if organ == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Require peeType")
	}

	if err := jawboneClient.eventIndigestion(time.Now(), organ); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error 500")
	}

	w.WriteHeader(http.StatusCreated)
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

	w.WriteHeader(http.StatusCreated)
}

// func handlerMealEvent(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
// 	mainItem := map[string][]string{}
// 	subItem := map[string][]string{}
// 	for key := range query {
// 		strList := query[key]
// 		if len(query[key]) == 0 {
// 			continue
// 		}

// 		if equalMealItemType("main", key) {
// 			subItem[key] = strList
// 		} else if equalMealItemType("sub", key) {
// 			mainItem[key] = strList
// 		}
// 	}
// 	 reqCreateMeal.
// 	//logrus.Debug()
// }

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
