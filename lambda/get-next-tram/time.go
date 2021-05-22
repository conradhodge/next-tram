package main

import "time"

var (
	// TimeNow monkey patches the time.Now function for unit-testing
	TimeNow = time.Now
)
