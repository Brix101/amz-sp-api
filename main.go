package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", endpoint+"/sellers/v1/marketplaceParticipations", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	err = sp.AuthorizeRequest(req)
	if err != nil {
		fmt.Println("Error authorizing request:", err)
	}
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	// Print the response body
	fmt.Println(string(body))
}
