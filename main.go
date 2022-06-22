package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Status struct {
	Water  int    `json:"water"`
	Wind   int    `json:"wind"`
	Status string `json:"status"`
}

var Weather Status

var PORT = ":8080"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/assignment-3", GetStatus)

	go GenereteToJSON()

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func GenereteToJSON() {
	for {
		Weather = Status{
			Water:  rand.Intn(100),
			Wind:   rand.Intn(100),
			Status: "aman",
		}

		if Weather.Water < 5 || Weather.Wind < 6 {
			Weather.Status = "aman"
		}
		if (Weather.Water >= 6 && Weather.Water <= 8) || (Weather.Wind >= 7 && Weather.Wind <= 15) {
			Weather.Status = "siaga"
		}
		if Weather.Water > 8 || Weather.Wind > 15 {
			Weather.Status = "bahaya"
		}

		//write json file
		jsonString, _ := json.Marshal(&Weather)
		ioutil.WriteFile("weather.json", jsonString, os.ModePerm)

		time.Sleep(15 * time.Second)
	}
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("weather.json")
	json.Unmarshal(file, &Weather)

	tpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	context := Status{
		Water:  Weather.Water,
		Wind:   Weather.Wind,
		Status: Weather.Status,
	}

	tpl.Execute(w, context)
}
