package main

import (
	"fmt"

	"net/url"

	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/schema"
	shortid "github.com/ventu-io/go-shortid"
)

const apiURL = "https://jawbone.com"

var (
	decoder = schema.NewDecoder()
	encoder = schema.NewEncoder()
)

type jawbone struct {
	ID    string
	Token *token
}

/*
note	string	Title / description of the meal. Used as both title and note.
sub_type	int	Meal type. 1=Breakfast, 2=Lunch, 3=Dinner, 4=Pre-Workout, 5=Post-Workout, 6=Snack
image_url	URI	URI of the meal image
photo	binary	Binary contents of the meal image
place_lat	float	Latitude of the location where the meal was created
place_lon	float	Longitude of the location where the meal was created
place_acc	float	Accuracy (meters) of the location where the meal was created
place_name	string	Name of the location where the meal was created
time_created	int	Epoch timestamp when the meal was created
tz	string	Time zone when this event was generated. Whenever possible, Olson format (e.g., "America/Los Angeles") will be returned, otherwise the GMT offset (e.g., "GMT+0800") will be returned.
share	boolean	Set whether to share event on user's public feed. Will not override if user's privacy setting is set to not share.
items	JSON-encoded list	See list below
*/

type reqCreateMeal struct {
	Note        string              `schema:"note"`
	SubType     int                 `schema:"sub_type"`
	ImageURL    string              `schema:"image_url"`
	Photo       string              `schema:"photo"`
	PlaceLat    float64             `schema:"place_lat"`
	PlaceLon    float64             `schema:"place_lon"`
	PlaceAcc    float64             `schema:"place_acc"`
	PlaceName   string              `schema:"place_name"`
	TimeCreated int                 `schema:"time_created"`
	TimeZone    string              `schema:"tz"`
	Share       bool                `schema:"share"`
	Items       []reqCreateMealItem `schema:"-"`
}

/*
Item property	Type	Description
name	string	Name of the meal item
description	string	Description of the meal item
amount	float	Amount of "measurement" (e.g. 100)
measurement	String	Unit of measurement (e.g. grams)
type	int	Quantity size by container. 1=Plate, 2=Cup, 3=Bowl, 4=Scale, 5=Glass
sub_type	int	Type of food. 1=Drink, 2=Food
food_categories	list	Used to log water. ["water"] = one glass of water
category	string	Category name (free text)
food_type	int	1=Generic, 2=Restaurant, 3=Brand, 4=Personal
calcium	int	Calcium (in milligrams)
calories	int	Calories
carbohydrate	float	Carbohydrate content (in grams)
cholesterol	float	Cholesterol content (in milligrams)
fiber	float	Fiber content (in grams)
protein	float	Protein content (in grams)
saturated_fat	float	Saturated fat content (in grams)
sodium	float	Sodium content (in milligrams)
sugar	float	Sugar content (in grams)
unsaturated_fat	float	Unsaturated fat content (in grams)
caffeine	float	Caffeine content (in milligrams)
*/
type reqCreateMealItem struct {
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	Amount         float64  `json:"amount,omitempty"`
	Measurement    string   `json:"measurement,omitempty"`
	Type           int      `json:"type,omitempty"`
	SubType        int      `json:"sub_type,omitempty"`
	FoodCategories []string `json:"food_categories,omitempty"`
	Category       string   `json:"category,omitempty"`
	FoodType       int      `json:"food_type,omitempty"`
	Calcium        int      `json:"calcium,omitempty"`
	Calories       int      `json:"calories,omitempty"`
	Carbohydrate   float64  `json:"carbohydrate,omitempty"`
	Cholesterol    float64  `json:"cholesterol,omitempty"`
	Fiber          float64  `json:"fiber,omitempty"`
	Protein        float64  `json:"protein,omitempty"`
	SaturatedFat   float64  `json:"saturated_fat,omitempty"`
	Sodium         float64  `json:"sodium,omitempty"`
	Sugar          float64  `json:"sugar,omitempty"`
	UnsaturatedFat float64  `json:"unsaturated_fat,omitempty"`
	Caffeine       float64  `json:"caffeine,omitempty"`
}

/*
Parameter	Type	Description
title	string	Name of the event (used in the feed story). 255 characters max.
verb	string	Verb to indicate user action (used in the feed story). 34 characters max.
attributes	json	Set of attributes associated with the event (for partner data only, not exposed in feed).
note	string	Description of the event. 512 characters max. URL links can be included to link outside the app. HTML formatting is not available.
image_url	URI	URI of the event's image
place_lat	float	Latitude of the location where the event was created
place_lon	float	Longitude of the location where the event was created
place_acc	float	Accuracy (meters) of the location where the event was created
place_name	string	Name of the location where the event was created
time_created	int	Unix timestamp when the event was created
tz	string	Time zone when this event was generated. Whenever possible, Olson format (e.g., "America/Los Angeles") will be returned, otherwise the GMT offset (e.g., "GMT+0800") will be returned.
share	boolean	Set whether to share event on user's public feed. Will not override if user's privacy setting is set to not share.
*/

type reqCreateCustom struct {
	Title       string                 `schema:"title"`
	Verb        string                 `schema:"verb"`
	Attributes  map[string]interface{} `schema:"-"`
	Note        string                 `schema:"note"`
	ImageURL    string                 `schema:"image_url"`
	PlaceLat    float64                `schema:"place_lat"`
	PlaceLon    float64                `schema:"place_lon"`
	PlaceAcc    float64                `schema:"place_acc"`
	PlaceName   string                 `schema:"place_name"`
	TimeCreated int                    `schema:"time_created"`
	TimeZone    string                 `schema:"tz"`
	Share       bool                   `schema:"share"`
}

type token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (cm reqCreateMeal) ConvForm() *url.Values {
	urlValue := url.Values{}
	if err := encoder.Encode(cm, urlValue); err != nil {
		logrus.WithError(err).Error("Faild Encode")
	}

	jsonItemBytes, err := json.Marshal(cm.Items)
	if err == nil {
		urlValue.Add("items", string(jsonItemBytes))
	} else {
		logrus.WithError(err).Error("Faild item marshal")
	}

	return &urlValue
}

func (cc reqCreateCustom) ConvForm() *url.Values {
	urlValue := url.Values{}
	encoder.Encode(cc, urlValue)

	jsonItemBytes, err := json.Marshal(cc.Attributes)
	if err == nil {
		urlValue.Add("attributes", string(jsonItemBytes))
	} else {
		logrus.WithError(err).Error("Faild item marshal")
	}

	return &urlValue
}

func (j *jawbone) save() error {
	if j.Token == nil {
		return errInvalidToken
	}
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		return err
	}
	if j.ID, err = sid.Generate(); err != nil {
		return err
	}

	return jDB.CreateJawbone(*j)
}

func (j jawbone) createMeal(createMeal reqCreateMeal) error {
	requestParam := httpRequestParams{
		Method: "POST",
		URL:    apiURL + "/nudge/api/v.1.1/users/@me/meals",
		Forms:  createMeal.ConvForm(),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", j.Token.AccessToken),
			"Accept":        "application/json",
			"Host":          host,
		},
	}

	res, err := requestParam.request()
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		logrus.WithError(err).Errorf("Faild. status code = %d. %s", res.StatusCode, string(res.Body))
		return errCreateMeal
	}

	return nil
}

func (j jawbone) createEvent(createCustom reqCreateCustom) error {
	requestParam := httpRequestParams{
		Method: "POST",
		URL:    apiURL + "/nudge/api/v.1.1/users/@me/generic_events",
		Forms:  createCustom.ConvForm(),
		Headers: map[string]string{
			"Authorization":     fmt.Sprintf("Bearer %s", j.Token.AccessToken),
			j.Token.AccessToken: "",
			"Accept":            "application/json",
			"Host":              host,
		},
	}

	res, err := requestParam.request()
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		return errCreateMeal
	}

	return nil
}
