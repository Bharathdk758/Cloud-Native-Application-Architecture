package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	geocodeURL   = "http://www.mapquestapi.com/geocoding/v1/reverse"
	staticMapURL = "https://www.mapquestapi.com/staticmap/v5/map"
	apiKey       = "Tqu58YkjqvoFsFstVLrb7jCf50cI0oxh"
	weatherAPI   = "https://api.openweathermap.org/data/2.5/weather"
	weatherKey   = "e8014a2e0f54c204692bb99b189296ee"
)

type Location struct {
	Street     string `json:"street"`
	City       string `json:"adminArea5"`
	County     string `json:"adminArea4"`
	State      string `json:"adminArea3"`
	Country    string `json:"adminArea1"`
	PostalCode string `json:"postalCode"`
}

type GeocodeResponse struct {
	Results []struct {
		ProvidedLocation struct {
			LatLng struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"latLng"`
		} `json:"providedLocation"`
		Locations []Location `json:"locations"`
	} `json:"results"`
}

type WeatherResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

var states = map[string]float64{
	"AL": 0.66, "AK": 1.21, "AZ": 1.06, "AR": 0.59,
	"CA": 2.36, "CO": 1.62, "CT": 1.23, "DE": 1.13,
	"DC": 2.71, "FL": 1.02, "GA": 0.84, "HI": 1.00,
	"ID": 0.36, "IL": 0.32, "IN": 0.23, "IA": 0.24,
	"KS": 0.25, "KY": 0.23, "LA": 0.27, "ME": 0.31,
	"MD": 0.51, "MA": 0.63, "MI": 0.25, "MN": 0.37,
	"MS": 0.20, "MO": 0.26, "MT": 0.39, "NE": 0.26,
	"NV": 0.46, "NH": 0.43, "NJ": 0.54, "NM": 0.28,
	"NY": 0.51, "NC": 0.29, "ND": 0.32, "OH": 0.24,
	"OK": 0.22, "OR": 0.53, "PA": 0.29, "RI": 0.43,
	"SC": 0.27, "SD": 0.28, "TN": 0.28, "TX": 0.29,
	"UT": 0.48, "VT": 0.36, "VA": 0.44, "WA": 0.58,
	"WV": 0.19, "WI": 0.30, "WY": 0.36,
}

func getAddress(latitude, longitude float64) (string, error) {
	url := fmt.Sprintf("%s?key=%s&location=%f,%f&outFormat=json&thumbMaps=false", geocodeURL, apiKey, latitude, longitude)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if len(data.Results) == 0 || len(data.Results[0].Locations) == 0 {
		return "", fmt.Errorf("no location found for the given coordinates")
	}

	location := data.Results[0].Locations[0]
	return fmt.Sprintf("%s, %s, %s, %s, %s %s", location.Street, location.City, location.County, location.State, location.Country, location.PostalCode), nil
}

func getState(latitude, longitude float64) (string, error) {
	url := fmt.Sprintf("%s?key=%s&location=%f,%f&outFormat=json&thumbMaps=false", geocodeURL, apiKey, latitude, longitude)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if len(data.Results) == 0 || len(data.Results[0].Locations) == 0 {
		return "", fmt.Errorf("no location found for the given coordinates")
	}

	location := data.Results[0].Locations[0]
	state := location.State
	if len(state) == 2 {
		state = strings.ToUpper(state)
	} else {
		state = strings.Title(strings.ToLower(state))
	}

	adjacent, ok := states[state]
	if !ok {
		return "", fmt.Errorf("state not found in the map")
	}

	return fmt.Sprintf("%.2f", adjacent), nil
}

func getWeather(latitude, longitude float64) (float64, error) {
	url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s", weatherAPI, latitude, longitude, weatherKey)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	if len(data.Weather) == 0 {
		return 0, fmt.Errorf("no weather information found")
	}

	return data.Main.Temp, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	latitude, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid latitude")
		return
	}
	longitude, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid longitude")
		return
	}

	address, err := getAddress(latitude, longitude)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error getting address:", err)
		return
	}
	adjacent, err := getState(latitude, longitude)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error getting adjacent value:", err)
		return
	}
	temperature, err := getWeather(latitude, longitude)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error getting weather information:", err)
		return
	}

	response := map[string]interface{}{
		"address":        address,
		"state_adjacent": adjacent,
		"temperature":    temperature,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server started on port 10000")
	http.ListenAndServe(":10000", nil)
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
}
