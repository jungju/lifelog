package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"

	"github.com/boltdb/bolt"
)

type jawbone struct {
	ID          string
	AccessToken string
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
	Note        string              `json:"note"`
	SubType     int                 `json:"sub_type"`
	ImageURL    string              `json:"image_url"`
	Photo       string              `json:"photo"`
	PlaceLat    float64             `json:"place_lat"`
	PlaceLon    float64             `json:"place_lon"`
	PlaceAcc    float64             `json:"place_acc"`
	PlaceName   string              `json:"place_name"`
	TimeCreated int                 `json:"time_created"`
	TimeZone    string              `json:"tz"`
	Share       bool                `json:"share"`
	Items       []reqCreateMealItem `json:"items"`
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
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Amount         float64  `json:"amount"`
	Measurement    string   `json:"measurement"`
	Type           int      `json:"type"`
	SubType        int      `json:"sub_type"`
	FoodCategories []string `json:"food_categories"`
	Category       string   `json:"category"`
	FoodType       int      `json:"food_type"`
	Calcium        int      `json:"calcium"`
	Calories       int      `json:"calories"`
	Carbohydrate   float64  `json:"carbohydrate"`
	Fiber          float64  `json:"fiber"`
	Protein        float64  `json:"protein"`
	SaturatedFat   float64  `json:"saturated_fat"`
	Sodium         float64  `json:"sodium"`
	Sugar          float64  `json:"sugar"`
	UnsaturatedFat float64  `json:"unsaturated_fat"`
	Caffeine       float64  `json:"caffeine"`
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
	Title       string                 `json:"title"`
	Verb        string                 `json:"verb"`
	Attributes  map[string]interface{} `json:"attributes"`
	Note        string                 `json:"note"`
	ImageURL    string                 `json:"image_url"`
	PlaceLat    float64                `json:"place_lat"`
	PlaceLon    float64                `json:"place_lon"`
	PlaceAcc    float64                `json:"place_acc"`
	PlaceName   string                 `json:"place_name"`
	TimeCreated int                    `json:"time_created"`
	TimeZone    string                 `json:"tz"`
	Share       bool                   `json:"share"`
}

func (j jawbone) save() error {
	j.ID = newUUID()
	return db.Update(func(tx *bolt.Tx) error {
		jawbones, err := tx.CreateBucketIfNotExists([]byte("jawbones"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		jawboneBytes, err := json.Marshal(&j)
		if err != nil {
			return nil
		}

		return jawbones.Put([]byte(j.ID), jawboneBytes)
	})
}

func (j jawbone) createMeal(createMeal reqCreateMeal) error {
	requestParam := httpRequestParams{
		Method:  "POST",
		URL:     "/nudge/api/v.1.1/users/@me/meals",
		Body:    &createMeal,
		Forms:   nil,
		Headers: map[string]string{"Authorization": "Bearer", j.AccessToken: ""},
	}

	res, err := requestParam.request()
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errCreateMeal
	}

	return nil
}

func (j jawbone) createEvent(createCustom reqCreateCustom) error {
	requestParam := httpRequestParams{
		Method:  "POST",
		URL:     "/nudge/api/v.1.1/users/@me/generic_events",
		Body:    &createCustom,
		Forms:   nil,
		Headers: map[string]string{"Authorization": "Bearer", j.AccessToken: ""},
	}

	res, err := requestParam.request()
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errCreateMeal
	}

	return nil
}

func (j jawbone) drinkWater(cups int) error {
	createWater := reqCreateMeal{
		Note:    "물먹기",
		SubType: 0,
		//ImageURL: "",
		Items: []reqCreateMealItem{
			reqCreateMealItem{
				Name:           "Water",
				Amount:         float64(cups),
				Measurement:    "cpus",
				Type:           2,
				FoodCategories: []string{"water"},
				FoodType:       1,
				Calories:       0,
			},
		},
	}
	return j.createMeal(createWater)
}

//dumpType 1: 매우좋은, 2:색깔이 안좋은, 3:물
func (j jawbone) takeDump(dumpType int, pain bool, constipation bool, blood bool) error {
	createDumpEvnet := reqCreateCustom{
		Title: "대변",
		Verb:  "보다",
		Attributes: map[string]interface{}{
			"dumpType":     dumpType,
			"pain":         pain,
			"constipation": constipation,
			"blood":        blood,
		},
		Note: fmt.Sprintf("상태 : %d, 고통: %t, 변비: %t, 피: %t", dumpType, pain, constipation, blood),
	}
	return j.createEvent(createDumpEvnet)
}

//peeType 1: 매우좋은, 2:색깔이 안좋은
func (j jawbone) pee(peeType int, blood bool) error {
	createDumpEvnet := reqCreateCustom{
		Title: "소변",
		Verb:  "보다",
		Attributes: map[string]interface{}{
			"peeType": peeType,
			"blood":   blood,
		},
		Note: fmt.Sprintf("상태 : %d, 피: %s", peeType, blood),
	}
	return j.createEvent(createDumpEvnet)
}

func getJawbone(id string) (*jawbone, error) {
	j := &jawbone{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("jawbones"))
		k := []byte(id)
		if err := json.Unmarshal(b.Get(k), j); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return j, nil
}

func newUUID() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "errorKey"
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
