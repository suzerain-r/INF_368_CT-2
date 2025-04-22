package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const apiBaseURL = "https://v6.exchangerate-api.com/v6"
const apiKey = "YOUR_API_KEY"

type ApiResponse struct {
	Result    string             `json:"result"`
	Rates     map[string]float64 `json:"conversion_rates"`
	ErrorType string             `json:"error-type"`
}

func main() {
	fmt.Print("From (e.g. USD): ")
	var fromCurrency string
	fmt.Scanln(&fromCurrency)

	fmt.Print("To (e.g. EUR): ")
	var toCurrency string
	fmt.Scanln(&toCurrency)

	fmt.Print("Amount: ")
	var amountStr string
	fmt.Scanln(&amountStr)

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Error: Invalid amount format.")
		return
	}

	url := fmt.Sprintf("%s/%s/latest/%s", apiBaseURL, apiKey, strings.ToUpper(fromCurrency))

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: Unable to reach API.")
		return
	}
	defer resp.Body.Close()

	var data ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error decoding API response.")
		return
	}

	if data.Result != "success" {
		fmt.Printf("API Error: %s\n", data.ErrorType)
		return
	}

	rate, exists := data.Rates[toCurrency]
	if !exists {
		fmt.Printf("Error: Invalid currency code '%s'\n", toCurrency)
		return
	}

	convertedAmount := amount * rate
	fmt.Printf("💱 %.2f %s = %.2f %s (Rate: %.4f)\n", amount, fromCurrency, convertedAmount, toCurrency, rate)
}
