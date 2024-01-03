package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	DB     *sql.DB
	Port   string
	Router *mux.Router
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func (app *App) Initialize() {
	DB, err := sql.Open("sqlite3", "../products.db")
	if err != nil {
		log.Fatal(err.Error())
	}

	app.DB = DB
	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/", welcome)
	app.Router.HandleFunc("/products", app.allProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.fetchProduct).Methods("GET")
	app.Router.HandleFunc("/products", app.newProduct).Methods("POST")

	app.Router.HandleFunc("/orders", app.allOrders).Methods("GET")
	app.Router.HandleFunc("/order/{id}", app.fetchOrder).Methods("GET")
	app.Router.HandleFunc("/orders", app.newOrder).Methods("POST")
	app.Router.HandleFunc("/orderItems", app.newOrderItems).Methods("POST")
}

func (app *App) allProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		fmt.Printf("getProducts error: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

func (app *App) fetchProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var p product
	p.ID, _ = strconv.Atoi(id)
	err := p.getProduct(app.DB)
	if err != nil {
		fmt.Printf("getProduct error: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (app *App) newProduct(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var p product
	json.Unmarshal(reqBody, &p)

	err := p.createProduct(app.DB)
	if err != nil {
		fmt.Printf("newProduct error: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (app *App) allOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := getOrders(app.DB)
	if err != nil {
		fmt.Printf("getOrders error: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, orders)
}

func (app *App) fetchOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var o order
	o.ID, _ = strconv.Atoi(id)
	err := o.getOrder(app.DB)
	if err != nil {
		fmt.Printf("getOrder error: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, o)
}

func (app *App) newOrder(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)

	var o order
	json.Unmarshal(reqBody, &o)
	err := o.createOrder(app.DB)
	if err != nil {
		fmt.Printf("newOrder error: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, item := range o.Items {
		oi := item
		oi.OrderID = o.ID
		err := oi.createOrderItem(app.DB)
		if err != nil {
			fmt.Printf("newOrderItem error: %s", err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, o)
	}
}

func (app *App) newOrderItems(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var ois []orderItem
	json.Unmarshal(reqBody, &ois)

	for _, item := range ois {
		oi := item
		err := oi.createOrderItem(app.DB)
		if err != nil {
			fmt.Printf("newOrderItem error: %s", err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	respondWithJSON(w, http.StatusOK, ois)
}

func (app *App) Run() {
	fmt.Println("Starting server on port ", app.Port)
	log.Fatal(http.ListenAndServe(app.Port, app.Router))
}

// Helper functions
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	//Convert payload to json
	response, _ := json.Marshal(payload)

	//Set header and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
