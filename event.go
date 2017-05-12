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

func (j jawbone) eventIcedAmericano(createdtime time.Time) error {
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

func (j jawbone) eventIcedLatte(createdtime time.Time) error {
	return j.createMeal(reqCreateMeal{
		Note:    "아이스 라떼",
		SubType: 6,
		//ImageURL:    "https://jawbone.com/ver/static/images/up/nutrition/categories/Food_icn_drink@2x_high.png",
		TimeZone:    "Asia/Seoul",
		TimeCreated: int(createdtime.Unix()),
		Items: []reqCreateMealItem{
			reqCreateMealItem{
				Name:         "아이스 라떼",
				Amount:       float64(1),
				Measurement:  "잔",
				Type:         5,
				SubType:      1,
				FoodType:     1,
				Category:     "coffee",
				Carbohydrate: 13,
				Protein:      7,
				Sodium:       110,
				Calories:     80,
				Sugar:        10,
			},
		},
	})
}

//dumpType 1: 매우좋은, 2:색깔이 안좋은, 3:물
func (j jawbone) eventPooh(createdtime time.Time, dumpType int, pain bool, constipation bool, blood bool) error {
	return j.createEvent(reqCreateCustom{
		Title:       "대변",
		Verb:        "배변",
		TimeZone:    "Asia/Seoul",
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

//peeType 1: 좋은, 2:나쁜
//color:1: 하얀, 2: 약한노란, 3: 진한노란
func (j jawbone) eventUrine(createdtime time.Time, peeType int, color int, blood bool) error {
	return j.createEvent(reqCreateCustom{
		Title:       "소변",
		Verb:        "누음",
		TimeZone:    "Asia/Seoul",
		TimeCreated: int(createdtime.Unix()),
		Attributes: map[string]interface{}{
			"peeType": peeType,
			"color":   color,
			"blood":   blood,
		},
		Note: fmt.Sprintf("상태 : %d, 피: %s", peeType, blood),
	})
}

//direction 1: 오른쪽, 2:왼쪽
func (j jawbone) eventMigraine(createdtime time.Time, direction int) error {
	return j.createEvent(reqCreateCustom{
		Title:       "편두통",
		Verb:        "아픔",
		TimeCreated: int(createdtime.Unix()),
		Attributes: map[string]interface{}{
			"direction": direction,
		},
		Note: fmt.Sprintf("부분 : %d", direction),
	})
}

//organ 1: 위, 2: 장
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
