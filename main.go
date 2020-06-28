package main

import (
	"github.com/tonydmorris/takeaway_payments/app"
	"github.com/tonydmorris/takeaway_payments/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
