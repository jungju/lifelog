package main

import (
	"fmt"
	"time"
)

func (j jawbone) eventWater(createdtime time.Time, cups int) error {
	return j.createMeal(reqCreateMeal{
		Note:    "Water",
		SubType: 4,
		//ImageURL: "",
		TimeZone:    "Asia/Seoul",
		Share:       true,
		TimeCreated: int(createdtime.Unix()),
		Items: []reqCreateMealItem{
			reqCreateMealItem{
				Amount:         float64(cups),
				Type:           5,
				SubType:        1,
				FoodCategories: []string{"water"},
			},
		},
	})
}

func (j jawbone) eventCoffee(createdtime time.Time) error {
	return j.createMeal(reqCreateMeal{
		Note:    "아이스 아메리카노",
		SubType: 6,
		//ImageURL:    "https://jawbone.com/ver/static/images/up/nutrition/categories/Food_icn_drink@2x_high.png",
		TimeZone:    "Asia/Seoul",
		TimeCreated: int(createdtime.Unix()),
		Items: []reqCreateMealItem{
			reqCreateMealItem{
				Name:         "아이스 아메리카노",
				Amount:       float64(1),
				Measurement:  "잔",
				Type:         5,
				SubType:      1,
				FoodType:     1,
				Category:     "coffee",
				Carbohydrate: 2,
				Protein:      1,
				Sodium:       5,
				Calories:     10,
			},
		},
	})
}

//dumpType 1: 매우좋은, 2:색깔이 안좋은, 3:물
func (j jawbone) eventPooh(createdtime time.Time, dumpType int, pain bool, constipation bool, blood bool) error {
	return j.createEvent(reqCreateCustom{
		Title:       "대변",
		Verb:        "보다",
		TimeCreated: int(createdtime.Unix()),
		Attributes: map[string]interface{}{
			"dumpType":     dumpType,
			"pain":         pain,
			"constipation": constipation,
			"blood":        blood,
		},
		Note: fmt.Sprintf("상태 : %d, 고통: %t, 변비: %t, 피: %t", dumpType, pain, constipation, blood),
	})
}

//peeType 1: 매우좋은, 2:색깔이 안좋은
func (j jawbone) eventUrine(createdtime time.Time, peeType int, blood bool) error {
	return j.createEvent(reqCreateCustom{
		Title:       "소변",
		Verb:        "누음",
		TimeCreated: int(createdtime.Unix()),
		Attributes: map[string]interface{}{
			"peeType": peeType,
			"blood":   blood,
		},
		Note: fmt.Sprintf("상태 : %d, 피: %s", peeType, blood),
	})
}

func (j jawbone) eventMigraine(createdtime time.Time, direction int) error {
	return j.createEvent(reqCreateCustom{
		Title:       "소변",
		Verb:        "누음",
		TimeCreated: int(createdtime.Unix()),
		Attributes: map[string]interface{}{
			"direction": direction,
		},
		Note: fmt.Sprintf("부분 : %d", direction),
	})
}

func (j jawbone) eventIndigestion(createdtime time.Time, organ int) error {
	return j.createEvent(reqCreateCustom{
		Title:       "소화불량",
		Verb:        "아픔",
		TimeCreated: int(createdtime.Unix()),
		Attributes: map[string]interface{}{
			"organ": organ,
		},
		Note: fmt.Sprintf("기관: %s", organ),
	})
}
