package handlers

import (
	_ "embed"
	"time"
)

//go:embed assets/logo.png
var logoBytes []byte

var timeNow = time.Now

const oneHour = time.Hour
