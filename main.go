package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Brix101/amz-sp-api/pkg/spsdk"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sp, err := spsdk.NewSellingPartner(&spsdk.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
	})

	if err != nil {
		panic(err)
	}

	endpoint := "https://sellingpartnerapi-na.amazon.com"
	queryUrl, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}

	basePath := fmt.Sprintf("/sellers/v1/marketplaceParticipations")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		panic(err)
	}

	log.Println(queryUrl)
	// Run the function immediately
	makeRequest(queryUrl.String(), sp)

	// Create a ticker to trigger the function every hour
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	// Run the function every time the ticker ticks
	for range ticker.C {
		makeRequest(queryUrl.String(), sp)
	}
}

func makeRequest(endpoint string, sp *spsdk.SellingPartner) {
	// Create a new GET request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Authorize the request (assuming sp.AuthorizeRequest exists)
	err = sp.AuthorizeRequest(req)
	if err != nil {
		fmt.Println("Error authorizing request:", err)
		return
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Print the response body
	log.Printf("Response received: %s at %s\n", string(body), time.Now().Format(time.RFC3339))
	log.Printf(time.Now().Format(time.RFC3339))
}
