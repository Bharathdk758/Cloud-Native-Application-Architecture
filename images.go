package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type PhotosResponse struct {
	Photos struct {
		Page    int `json:"page"`
		Pages   int `json:"pages"`
		PerPage int `json:"perpage"`
		Total   int `json:"total"`
		Photo   []struct {
			ID       string `json:"id"`
			Secret   string `json:"secret"`
			Server   string `json:"server"`
			Farm     int    `json:"farm"`
			Title    string `json:"title"`
			IsPublic int    `json:"ispublic"`
			IsFriend int    `json:"isfriend"`
			IsFamily int    `json:"isfamily"`
		} `json:"photo"`
	} `json:"photos"`
	Stat string `json:"stat"`
}

func main() {
	// Set the Flickr API endpoint and parameters
	endpoint := "https://api.flickr.com/services/rest/"
	params := url.Values{}
	params.Set("method", "flickr.photos.search")
	params.Set("api_key", "1ab482ed35ec3c5046bae21229994de2")
	params.Set("format", "json")
	params.Set("nojsoncallback", "1")
	params.Set("per_page", "8") // Set per_page to 4

	// Create an HTTP handler function to serve requests for /photos.json
	http.HandleFunc("/photos.json", func(w http.ResponseWriter, r *http.Request) {
		// Set the CORS headers
		setCORSHeaders(w)

		params.Set("lat", r.URL.Query().Get("lat"))
		params.Set("lon", r.URL.Query().Get("lon"))

		// Build the URL and send the request
		reqURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
		resp, err := http.Get(reqURL)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to fetch photos", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Parse the response
		var photos PhotosResponse
		err = json.NewDecoder(resp.Body).Decode(&photos)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to decode photos", http.StatusInternalServerError)
			return
		}

		// Write the photos variable to the HTTP response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(photos)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to encode photos", http.StatusInternalServerError)
			return
		}
	})

	// Start the HTTP server
	fmt.Println("Listening on http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
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
