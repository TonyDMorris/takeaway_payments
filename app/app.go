package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tonydmorris/takeaway_payments/app/handler"
	"github.com/tonydmorris/takeaway_payments/config"
)

// App has router and db instances
type App struct {
	Router      *mux.Router
	DB          *sql.DB
	HTTPClient  *http.Client
	StrapiToken string
	StapiURL    string
}
type Credentials struct {
	Jwt string `json:jwt`
}

type Authentication struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {

	// initalise db connection
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	config.DBHost, config.DBPort, config.DBUsername, config.DBPassword, config.DBName)
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	// a.DB = db

	// authenticate payment service with strapi
	a.StapiURL = config.StrapiURL
	authURL := fmt.Sprintf("%v/auth/local", a.StapiURL)
	var authentication Authentication
	authentication.Identifier = config.ServiceIdentifier
	authentication.Password = config.ServicePassword
	body, err := json.Marshal(authentication)

	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	a.HTTPClient = &http.Client{Timeout: time.Second * 10, Transport: &http.Transport{
		DisableKeepAlives:   false,
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 0,
		IdleConnTimeout:     time.Second * 5,
	}}
	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var credentials Credentials
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &credentials)
	if err != nil {
		panic(err)
	}

	a.StrapiToken = credentials.Jwt
	a.Router = mux.NewRouter()
	a.setRouters()
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// Routing for handling the projects

	a.Post("/payments", a.handleRequest(handler.HandlePayments))

}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(client *http.Client, strapiUrl string, strapiToken string, db *sql.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.HTTPClient, a.StapiURL, a.StrapiToken, a.DB, w, r)
	}
}
