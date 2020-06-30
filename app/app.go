package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tonydmorris/takeaway_payments/config"
)

// App has router and db instances
type App struct {
	Router      *mux.Router
	DB          *sql.DB
	StrapiToken string
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {

	// initalise db connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUsername, config.DBPassword, config.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	a.DB = db

	// authenticate payment service with strapi
	authURL := fmt.Sprintf("%v/auth/local", config.StrapiURL)
	authentication := make(map[string]string)
	authentication["identifier"] = config.ServiceIdentifier
	authentication["password"] = config.ServicePassword
	body, err := json.Marshal(authentication)
	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	credentials := make(map[string]string)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bodyBytes, &credentials)
	if err != nil {
		panic(err)
	}
	a.StrapiToken = credentials["jwt"]
	a.Router = mux.NewRouter()
	a.setRouters()
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// Routing for handling the projects

	a.Post("/projects", a.handleRequest())

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

type RequestHandlerFunction func(db *sql.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}
