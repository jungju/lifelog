package main

import "errors"

var (
	errCreateMeal   = errors.New("ErrCreateMeal")
	errInvalidToken = errors.New("ErrInvalidToken")
)
