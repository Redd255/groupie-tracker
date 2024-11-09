package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var tmp2 = template.Must(template.ParseFiles("templates/error.html"))

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
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
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		data := map[string]any{"code": http.StatusInternalServerError, "msg": "error parsing file "}
		w.WriteHeader(http.StatusInternalServerError)
		tmp2.Execute(w, data)
		return
	}

	data, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		data := map[string]any{"code": http.StatusInternalServerError, "msg": "failed fetshing data"}
		w.WriteHeader(http.StatusInternalServerError)
		tmp2.Execute(w, data)
		return
	}
	defer data.Body.Close()

	var artists []Artist
	if err := json.NewDecoder(data.Body).Decode(&artists); err != nil {
		data := map[string]any{"code": http.StatusInternalServerError, "msg": "failed to proccess api data"}
		w.WriteHeader(http.StatusInternalServerError)
		tmp2.Execute(w, data)
		return
	}

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
	tmpl, err := template.ParseFiles("templates/secondpage.html")
	if err != nil {
		data := map[string]any{"code": http.StatusInternalServerError, "msg": "error parsing file "}
		w.WriteHeader(http.StatusInternalServerError)
		tmp2.Execute(w, data)
		return
	}

	idStr := r.URL.Path[len("/details/"):]

	if idStr < "0" || idStr > "52" {
		data := map[string]any{"code": http.StatusNotFound, "msg": "page not found"}
		w.WriteHeader(http.StatusNotFound)
		tmp2.Execute(w, data)
		return
	}

	url1 := "https://groupietrackers.herokuapp.com/api/artists/" + idStr
	data1, err := http.Get(url1)
	if err != nil {
		data := map[string]any{"code": http.StatusInternalServerError, "msg": "failed fetshing data"}
		w.WriteHeader(http.StatusInternalServerError)
		tmp2.Execute(w, data)
		return
	}
	defer data1.Body.Close()

	var artist Artist
	if err := json.NewDecoder(data1.Body).Decode(&artist); err != nil {
		data := map[string]any{"code": http.StatusInternalServerError, "msg": "error decoding"}
		w.WriteHeader(http.StatusInternalServerError)
		tmp2.Execute(w, data)
		return
	}

	url2 := "https://groupietrackers.herokuapp.com/api/locations/" + idStr
	data2, err := http.Get(url2)
	if err != nil {
		data := map[string]any{"code": http.StatusInternalServerError, "msg": "error fetching locations"}
		w.WriteHeader(http.StatusInternalServerError)
		tmp2.Execute(w, data)
		return
	}
	defer data2.Body.Close()

	var locationsResponse LocationsResponse
	if err := json.NewDecoder(data2.Body).Decode(&locationsResponse); err != nil {
		data := map[string]any{"code": http.StatusInternalServerError, "msg": "errorr decoding"}
		w.WriteHeader(http.StatusInternalServerError)
		tmp2.Execute(w, data)
		return
	}

	locationsResponse.ID = artist.ID

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

	if idStr < "0" || idStr > "52" {
		data := map[string]any{"code": http.StatusNotFound, "msg": "page not found"}
		w.WriteHeader(http.StatusNotFound)
		tmp2.Execute(w, data)
		return
	}

	url := "https://groupietrackers.herokuapp.com/api/relation/" + idStr
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching relations data:", err)
		http.Error(w, "Failed to fetch relations data", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var relations Relations
	if err := json.NewDecoder(response.Body).Decode(&relations); err != nil {
		fmt.Println("Error decoding relations JSON:", err)
		http.Error(w, "Failed to decode relations data", http.StatusInternalServerError)
		return
	}

	url2 := "https://groupietrackers.herokuapp.com/api/dates/" + idStr
	response2, err := http.Get(url2)
	if err != nil {
		fmt.Println("Error fetching dates data:", err)
		http.Error(w, "Failed to fetch dates data", http.StatusInternalServerError)
		return
	}
	defer response2.Body.Close()

	var dates Dates
	if err := json.NewDecoder(response2.Body).Decode(&dates); err != nil {
		fmt.Println("Error decoding dates JSON:", err)
		http.Error(w, "Failed to decode dates data", http.StatusInternalServerError)
		return
	}

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
