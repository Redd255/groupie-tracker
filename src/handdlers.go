package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var tmp2 = template.Must(template.ParseFiles("templates/error.html"))

func HomePage(w http.ResponseWriter, r *http.Request) {
	// Check if the requested URL path starts with "/""
	if r.URL.Path != "/" {
		data := map[string]any{"code": http.StatusNotFound, "msg": "not found"}
		w.WriteHeader(http.StatusNotFound)
		tmp2.Execute(w, data)
		return
	}
	// Handle non-GET methods
	if r.Method != "GET" {
		data := map[string]any{"code": http.StatusMethodNotAllowed, "msg": "method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		tmp2.Execute(w, data)
		return
	}
	// Parse the artist page
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	url := "https://groupietrackers.herokuapp.com/api/artists"
	data := Fetch(url, w)

	var artists []Artist
	DecodeByUs(data, &artists, w)

	if err := tmpl.Execute(w, artists); err != nil {
		data := map[string]any{"code": http.StatusInternalServerError, "msg": "failed executing"}
		w.WriteHeader(http.StatusInternalServerError)
		tmp2.Execute(w, data)
		return
	}
}

func SecondPage(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/details/") {
		data := map[string]any{"code": http.StatusNotFound, "msg": "not found"}
		w.WriteHeader(http.StatusNotFound)
		tmp2.Execute(w, data)
		return
	}
	if r.Method != "GET" {
		data := map[string]any{"code": http.StatusMethodNotAllowed, "msg": "method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		tmp2.Execute(w, data)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/secondpage.html"))

	// Check the path id
	idStr := r.URL.Path[len("/details/"):]
	ids, _ := strconv.Atoi(idStr)
	if ids <= 0 || ids > 52 {
		data := map[string]any{"code": http.StatusNotFound, "msg": "page not found"}
		w.WriteHeader(http.StatusNotFound)
		tmp2.Execute(w, data)
		return
	}

	// Fetch artist data from the API
	url1 := "https://groupietrackers.herokuapp.com/api/artists/" + idStr
	data1 := Fetch(url1, w)

	// Decode the JSON response into the artist struct.
	var artist Artist
	DecodeByUs(data1, &artist, w)

	url2 := "https://groupietrackers.herokuapp.com/api/locations/" + idStr
	data2 := Fetch(url2, w)

	var locationsResponse LocationsResponse
	DecodeByUs(data2, &locationsResponse, w)

	pageData := SecondPageData{
		ID:        locationsResponse.ID,
		Artist:    artist,
		Locations: locationsResponse.Locations,
	}

	if err := tmpl.Execute(w, pageData); err != nil {
		data := map[string]any{"code": http.StatusInternalServerError, "msg": "error executing"}
		w.WriteHeader(http.StatusInternalServerError)
		tmp2.Execute(w, data)
		return
	}
}

func LastPage(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/lastpage/") {
		data := map[string]any{"code": http.StatusNotFound, "msg": "not found"}
		w.WriteHeader(http.StatusNotFound)
		tmp2.Execute(w, data)
		return
	}
	if r.Method != "GET" {
		data := map[string]any{"code": http.StatusMethodNotAllowed, "msg": "method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		tmp2.Execute(w, data)
		return
	}
	tmpl, err := template.ParseFiles("templates/thirdpage.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Path[len("/lastpage/"):]
	ids, _ := strconv.Atoi(idStr)
	if ids <= 0 || ids > 52 {
		data := map[string]any{"code": http.StatusNotFound, "msg": "page not found"}
		w.WriteHeader(http.StatusNotFound)
		tmp2.Execute(w, data)
		return
	}

	url := "https://groupietrackers.herokuapp.com/api/relation/" + idStr
	response := Fetch(url, w)

	var relations Relations
	DecodeByUs(response, &relations, w)

	url2 := "https://groupietrackers.herokuapp.com/api/dates/" + idStr
	response2 := Fetch(url2, w)

	var dates Dates
	DecodeByUs(response2, &dates, w)

	pageData := ThirdPageData{
		ID:             relations.ID,
		Dates:          dates.DATES,
		DatesLocations: relations.DATESLOCAT,
	}

	if err := tmpl.Execute(w, pageData); err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
