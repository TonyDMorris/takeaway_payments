package handler

type order struct {
	FoodItems []FoodItem
	DealItems []DealItem
	Address   Address
	Token     string // stripe payment intent token
}
type FoodItem struct {
	Id int
}
type DealItem struct {
	Id int
}

type Address struct {
	Postcode               string
	Phone                  string
	Number                 string // house or flat number
	Line1                  string
	Line2                  string
	AdditionalInstructions string // extra instructions for restaurant e.g. deliver around the back
}
