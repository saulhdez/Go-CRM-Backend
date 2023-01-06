package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid" // Required to generate an UUID for each customer
	"github.com/gorilla/mux" // Required for the HTTP method-based router
)

// Customer struct with JSON tags to properly encode/decode the JSON payloads
type Customer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     uint32 `json:"phone"`
	Contacted bool   `json:"contacted"`
}

// customers map which acts as the CRM's database
var customers = make(map[string]Customer)

/*
- HTTP handler: index
- Endpoint: "/"
- HTTP method: GET
- Response Content-Type: HTML
- Description: Serves an static HTML welcome page in the root route ("/").
- Response content: The static HTML welcome page.
*/
func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, "./static/index.html")
}

/*
  - Function: getCustomersSlice
  - Description: Utility function to convert the customers map into a slice of
    Customer structs so when retrieving all the customers they are correctly
    formatted in the JSON response.
  - Return value(s): A slice of Customer structs containing all the customers
    of the CRM's database.
*/
func getCustomersSlice() []Customer {
	customersSlice := make([]Customer, 0, len(customers))

	for _, customer := range customers {
		customersSlice = append(customersSlice, customer)
	}

	return customersSlice
}

/*
  - HTTP handler: getCustomers
  - Endpoint: "/customers"
  - HTTP method: GET
  - Response Content-Type: JSON
  - Description: Retrieves and returns all the customers already registered
    in the CRM.
  - Response content: A JSON array with all the customer objects (structs)
    contained in the CRM's database.
*/
func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getCustomersSlice())
}

/*
  - HTTP handler: getCustomer
  - Endpoint: "/customers/{id}"
  - HTTP method: GET
  - Response Content-Type: JSON
  - Description: Retrieves the data of a single customer based on the
    customer ID passed as a URL paramater.
  - Response content: A JSON object with the information of the requested
    customer if its ID was found in the database, or an empty JSON object
    literal in case the provided customer ID does not exists in the database.
*/
func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Retrieve the customer ID provided by the user in the URL params
	id := mux.Vars(r)["id"]

	if customer, exists := customers[id]; exists {
		// If the customer ID exists return its data with a 200 status code
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customer)
	} else {
		// Otherwise, return an empty JSON object literal and a 404 status code
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

/*
  - HTTP handler: addCustomer
  - Endpoint: "/customers"
  - HTTP method: POST
  - Response Content-Type: JSON
  - Description: Registers a new customer in the CRM's database using the
    customer data passed by the user in the JSON payload of the request.
  - Response content: A JSON object with the data of the newly created
    customer.
*/
func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Read, parse and store the JSON request data
	var newCustomer Customer
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newCustomer)

	// Assing an UUID to the new customer and add it to the CRM's database
	newCustomer.ID = uuid.New().String()
	customers[newCustomer.ID] = newCustomer
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newCustomer)
}

/*
  - HTTP handler: updateCustomer
  - Endpoint: "/customers/{id}"
  - HTTP method: PUT
  - Response Content-Type: JSON
  - Description: Updates the data of a single customer based on the customer ID
    passed as a URL paramater. For this process, the user passes in the request
    payload a JSON object containing just the customer fields to be updated
    with their new values. All fields not present in that JSON object are
    ignored and their values remain unchanged in the CRM's database. As a note,
    the customer's ID is never updated, even if the user passes a new value
    for it.
  - Response content: A JSON object containing all the data of the updated
    customer in case the passed customer ID was found in the database, or an
    empty JSON object literal in case the provided customer ID does not exists
    in the database.
*/
func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Retrieve the customer ID provided by the user in the URL params
	id := mux.Vars(r)["id"]

	if customer, exists := customers[id]; exists {
		// If the customer ID exists, we create two variables, receivedFields
		// which stores the fields that the user wants to update for the
		// customer; and the Customer struct variable updateCustomerData, which
		// stores all the new values. The values for the fields not provided by
		// the user get assiged their default empty/zero value in the
		// updateCustomerData variable when populating it.
		var receivedFields map[string]string
		var updateCustomerData Customer

		// Read, parse and store the JSON request data into the variables
		reqBody, _ := io.ReadAll(r.Body)
		json.Unmarshal(reqBody, &receivedFields)
		json.Unmarshal(reqBody, &updateCustomerData)

		// Iterate just over the fields that the user wants to update for the
		// customer and assign their new values
		for field := range receivedFields {
			switch field {
			case "name":
				customer.Name = updateCustomerData.Name
			case "role":
				customer.Role = updateCustomerData.Role
			case "email":
				customer.Email = updateCustomerData.Email
			case "phone":
				customer.Phone = updateCustomerData.Phone
			case "contacted":
				customer.Contacted = updateCustomerData.Contacted
			}
		}

		// Set the new customer data (struct) in the database (customers map)
		customers[id] = customer
		// Return a 200 code together with the newly updated customer data
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customer)
	} else {
		// If the customer ID was not found, then return a 404 status code
		// along with an empty JSON object literal
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

/*
  - HTTP handler: deleteCustomer
  - Endpoint: "/customers/{id}"
  - HTTP method: DELETE
  - Response Content-Type: JSON
  - Description: Deletes a single customer from the CRM's database based on the
    customer ID passed as a URL paramater.
  - Response content: If the provided customer ID was found in the database, the
    handler returns a JSON array with all the customer objects (structs) of the
    CRM's database, reflecting the deletion of the requested customer. If the
    provided customer ID is not found, then the response is an empty JSON
    object literal.
*/
func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Retrieve the customer ID provided by the user in the URL params
	id := mux.Vars(r)["id"]

	if _, exists := customers[id]; exists {
		// If the customer id is found, delete the customer from the CRM's
		// database and return a 200 status code along with a slice with all
		// the customer structs of the database
		delete(customers, id)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(getCustomersSlice())
	} else {
		// Otherwise, return an empty JSON object literal and a 404 status code
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

/*
  - Function: createSeedCustomers
  - Description: Utility function to seed the CRM's database with three initial
    customers.
  - Return value(s): None.
*/
func createSeedCustomers() {
	// Create the UUIDs for the new customers
	id1 := uuid.New().String()
	id2 := uuid.New().String()
	id3 := uuid.New().String()

	// Initialize the CRM's database with three new customers
	customers = map[string]Customer{
		id1: {id1, "Carlos", "Software Engineer", "carlos@example.com", 5553745, true},
		id2: {id2, "Samantha", "Manager", "samantha@example.com", 5550234, false},
		id3: {id3, "Santosh", "Director", "santosh@example.com", 5558721, false},
	}
}

func main() {
	// Initialize the CRM's database
	createSeedCustomers()

	// Create the gorilla/mux router to be used for the HTTP method-based routing
	router := mux.NewRouter().StrictSlash(true)
	port := "3000"

	// Register the HTTP handlers to their respective endpoints/routes and to
	// their assigned HTTP methods
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	// Start the server in the specified port and assign our custom router to
	// handle the incoming requests
	fmt.Println("Server started on port " + port + "...")
	http.ListenAndServe(":"+port, router)
}
