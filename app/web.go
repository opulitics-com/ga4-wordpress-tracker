package main

import (
	"encoding/json"
	"ga4-wordpresss-tracker/config"
	"ga4-wordpresss-tracker/service"
	"html/template"
	"log"
	"net/http"
)

func initWebServer() {
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatalf("Error loading template: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, "")
		if err != nil {
			http.Error(w, "Unable to render template", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		reports := service.GetReport()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(reports)
		w.WriteHeader(http.StatusOK)
	})
	log.Println("Server started at :" + config.AppConfig.Port)
	log.Fatal(http.ListenAndServe(":"+config.AppConfig.Port, nil))
}
