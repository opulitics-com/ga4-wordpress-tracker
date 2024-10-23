package main

import (
	"ga4-wordpresss-tracker/config"
)

func main() {
	config.LoadConfig()
	initWebServer()
}
