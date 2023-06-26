package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// Inputs and coefficients
	inputs := map[string]string{
		"Area":            "875.9677646745798",
		"Bedrooms":        "0.6686444174492727",
		"Bathrooms":       "0.32628884243475254",
		"Stories":         "0.4753911479186575",
		"Mainroad":        "0.18743314795398144",
		"Guestroom":       "0.05480303029780856",
		"Basement":        "0.09426090639781622",
		"Hotwaterheating": "0.01666196407928698",
		"Airconditioning": "0.10854823639827131",
		"Parking":         "0.17152131541760351",
		"Prefarea":        "0.06782700294669801",
		"Theta0":          "0.20485650505752898",
	}

	// Define handler function to handle HTTP requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse user inputs from HTTP request parameters

		area, _ := strconv.ParseFloat(r.FormValue("area"), 64)
		bedrooms, _ := strconv.Atoi(r.FormValue("bedrooms"))
		bathrooms, _ := strconv.Atoi(r.FormValue("bathrooms"))
		stories, _ := strconv.Atoi(r.FormValue("stories"))
		mainroad, _ := strconv.Atoi(r.FormValue("mainroad"))
		guestroom, _ := strconv.Atoi(r.FormValue("guestroom"))
		basement, _ := strconv.Atoi(r.FormValue("basement"))
		hotwaterheating, _ := strconv.Atoi(r.FormValue("hotwaterheating"))
		airconditioning, _ := strconv.Atoi(r.FormValue("airconditioning"))
		parking, _ := strconv.Atoi(r.FormValue("parking"))
		prefarea, _ := strconv.Atoi(r.FormValue("prefarea"))

		// Convert string inputs to float64 coefficients
		coefficients := make(map[string]float64, len(inputs))
		for key, value := range inputs {
			coefficient, _ := strconv.ParseFloat(value, 64)
			coefficients[key] = coefficient
		}

		// Calculate predicted value
		predictedValue := coefficients["Theta0"] +
			area*coefficients["Area"] +
			float64(bedrooms)*coefficients["Bedrooms"] +
			float64(bathrooms)*coefficients["Bathrooms"] +
			float64(stories)*coefficients["Stories"] +
			float64(mainroad)*coefficients["Mainroad"] +
			float64(guestroom)*coefficients["Guestroom"] +
			float64(basement)*coefficients["Basement"] +
			float64(hotwaterheating)*coefficients["Hotwaterheating"] +
			float64(airconditioning)*coefficients["Airconditioning"] +
			float64(parking)*coefficients["Parking"] +
			float64(prefarea)*coefficients["Prefarea"]

		var predictedValue2 int = int(predictedValue)
		// Output predicted value to HTTP response
		setCORSHeaders(w)
		response := map[string]int{"predictedValue": predictedValue2}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
		fmt.Println("Predicted value:", predictedValue)
	})

	// Start the HTTP server
	fmt.Println("Listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Function to set the CORS headers
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
}

//http://localhost:8000/?area=1000&bedrooms=2&bathrooms=1&stories=2&mainroad=1&guestroom=0&basement=1&hotwaterheating=0&airconditioning=1&parking=1&prefarea=1
