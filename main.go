package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof" // Import pprof for profiling
	"runtime"
	"time"
)

type Product struct {
	Code    string `json:"code"`
	Product struct {
		ID              string                 `json:"_id"`
		ProductName     string                 `json:"product_name"`
		Brands          string                 `json:"brands"`
		IngredientsText string                 `json:"ingredients_text"`
		Categories      string                 `json:"categories"`
		Traces          string                 `json:"traces"`
		Allergens       string                 `json:"allergens"`
		NutriscoreGrade string                 `json:"nutriscore_grade"`
		Nutriments      map[string]interface{} `json:"nutriments"`
	} `json:"product"`
}

// downloadJSON downloads the JSON data from the specified URL
func downloadJSON(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func serializeJSON(data Product, store *[]string) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error serializing JSON:", err)
		return
	}

	// Add the serialized JSON initially to the slice
	*store = append(*store, string(jsonData))

	initSleep := 5000 * time.Millisecond
	for {
		time.Sleep(initSleep)
		// Duplicate the slice content
		*store = append(*store, *store...)
		fmt.Println("Slice duplicated, current size:", len(*store))
		// Decrease the wait time
		initSleep += 500 * time.Millisecond
	}
}

func main() {
	fmt.Println("Starting exponential memory leak...")

	// Download the JSON file from the URL
	url := "https://world.openfoodfacts.org/api/v0/product/737628064502.json"
	jsonData, err := downloadJSON(url)
	if err != nil {
		log.Fatal("Error downloading JSON:", err)
	}

	// Deserialize JSON into the Product struct
	var data Product
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

	// Slice to store copies of the serialized JSON
	store := make([]string, 0)

	// Run the memory-leaking function in a goroutine
	go serializeJSON(data, &store)

	// Start pprof server on port 8080
	go func() {
		log.Println("Starting pprof server on http://localhost:8080/debug/pprof/")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("pprof server failed: %v", err)
		}
	}()

	// Monitor memory usage
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("Memory in use: %d MB\n", m.Alloc/1024/1024)
		time.Sleep(1 * time.Second)
	}
}
