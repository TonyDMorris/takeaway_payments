package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/stripe/stripe-go"
	stripeClient "github.com/stripe/stripe-go/client"
)

// CustomerOrder represents the incoming payment request from the fontend
type CustomerOrder struct {
	UUID         string     `json:"uuid"`
	RestaurantID string     `json:"restaurant_id"`
	FoodItems    []FoodItem `json:"food_items"`
	DealItems    []DealItem `json:"deal_items"`
	Address      Address    `json:"address"`
	PaymentToken string     `json:"payment_token"` // stripe payment intent token
}

// FoodItem represents a food item sent by the customer to be charged
type FoodItem struct {
	ID  int `json:"id"`
	Qty int `json:"qty"`
}

//DealItem represents a deal sent by the customer to be charged
type DealItem struct {
	ID  int `json:"id"`
	Qty int `json:"qty"`
}

//Address is the customers address for the order to be delivered too
type Address struct {
	Postcode               string
	Phone                  string
	Number                 string // house or flat number
	Line1                  string
	Line2                  string
	AdditionalInstructions string // extra instructions for restaurant e.g. deliver around the back
}

type Restaurant struct {
	ID         int    `json:"id"` //// decode restaurant for id
	Name       string `json:"name"`
	Credential struct {
		PublishableKey     string `json:"publishable_key"`
		PrivateKey         string `json:"private_key"`
		TestPublishableKey string `json:"test_publishable_key"`
		TestPrivateKey     string `json:"test_private_key"`
	} `json:"payment_credential"`
}

type CheckoutData struct {
	ClientSecret string `json:"client_secret"`
}

func calculateTotalPriceInPence(customerOrder CustomerOrder, restaurant Restaurant) int64 {
	return 1000
}

// HandlePayments is the entrypoint for a payment request
func HandlePayments(client *http.Client, strapiURL string, strapiToken string, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var customerOrder CustomerOrder
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &customerOrder)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	url := fmt.Sprintf("%v/restaurants/payment/%v", strapiURL, customerOrder.RestaurantID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+strapiToken)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var strapiResponse []Restaurant

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &strapiResponse)
	if err != nil {
		panic(err)
	}
	restaurant := strapiResponse[0]

	total := calculateTotalPriceInPence(customerOrder, restaurant)
	sc := stripeClient.New(restaurant.Credential.TestPrivateKey, nil)

	params := &stripe.PaymentIntentParams{
		Amount:              stripe.Int64(total),
		Currency:            stripe.String(string(stripe.CurrencyGBP)),
		SetupFutureUsage:    stripe.String(string(stripe.PaymentIntentSetupFutureUsageOffSession)),
		StatementDescriptor: stripe.String(restaurant.Name[:21]),
	}
	params.AddMetadata("order_id", customerOrder.UUID)
	pi, _ := sc.PaymentIntents.New(params)

	customerResponse := CheckoutData{
		ClientSecret: pi.ClientSecret,
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customerResponse)

}
