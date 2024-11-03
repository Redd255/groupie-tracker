package server

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	data, _ := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	var artists []Artist
	json.NewDecoder(data.Body).Decode(&artists)
	tmpl.Execute(w, artists)
}

func SecondPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/secondpage.html")
	idStr := r.URL.Path[len("/details/"):]

	url1 := "https://groupietrackers.herokuapp.com/api/artists/" + idStr
	data1, _ := http.Get(url1)
	var artist Artist
	json.NewDecoder(data1.Body).Decode(&artist)

	url2 := "https://groupietrackers.herokuapp.com/api/locations/" + idStr
	data2, _ := http.Get(url2)

	var locationsResponse LocationsResponse
	json.NewDecoder(data2.Body).Decode(&locationsResponse)

	locationsResponse.ID = artist.ID

	pageData := struct {
		ID        int
		Artist    Artist
		Locations []string
	}{
		ID:        locationsResponse.ID,
		Artist:    artist,
		Locations: locationsResponse.Locations,
	}

	tmpl.Execute(w, pageData)
}

func LastPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/thirdpage.html")

	idStr := r.URL.Path[len("/lastpage/"):]

	url := "https://groupietrackers.herokuapp.com/api/relation/" + idStr
	response, _ := http.Get(url)

	var relations Relations
	json.NewDecoder(response.Body).Decode(&relations)

	url2 := "https://groupietrackers.herokuapp.com/api/dates/" + idStr
	response2, _ := http.Get(url2)

	var dates Dates
	json.NewDecoder(response2.Body).Decode(&dates)

	pageData := struct {
		ID             int                 `json:"id"`
		Dates          []string            `json:"dates"`
		DatesLocations map[string][]string `json:"datesLocations"`
	}{
		ID:             relations.ID,
		Dates:          dates.DATES,
		DatesLocations: relations.DATESLOCAT,
	}

	tmpl.Execute(w, pageData)
}
