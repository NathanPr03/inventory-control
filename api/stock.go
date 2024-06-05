package handler

import (
	"encoding/json"
	"github.com/NathanPr03/price-control/pkg/db"
	"inventory-control/api/generated"
	"net/http"
)

var dbConnection, _ = db.ConnectToDb()

func StockHandler(w http.ResponseWriter, r *http.Request) {
	query := "SELECT name FROM products WHERE remaining_stock < 10"
	rows, err := dbConnection.Query(query)
	if err != nil {
		http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []string
	for rows.Next() {
		var productName string
		if err := rows.Scan(&productName); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, productName)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error with rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string][]string{"products": products}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func IncrementStock(w http.ResponseWriter, r *http.Request) {
	productName := r.URL.Query().Get("productName")
	if productName == "" {
		http.Error(w, "Missing productName parameter", http.StatusBadRequest)
		return
	}

	query := "UPDATE products SET remaining_stock = remaining_stock + 1 WHERE name = $1"
	_, err := dbConnection.Exec(query, productName)
	if err != nil {
		http.Error(w, "Error updating stock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Stock incremented successfully"})
}

func DecrementStock(w http.ResponseWriter, r *http.Request) {
	productName := r.URL.Query().Get("productName")
	if productName == "" {
		http.Error(w, "Missing productName parameter", http.StatusBadRequest)
		return
	}

	query := "UPDATE products SET remaining_stock = remaining_stock - 1 WHERE name = $1"
	_, err := dbConnection.Exec(query, productName)
	if err != nil {
		http.Error(w, "Error updating stock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Stock decremented successfully"})
}

func ChangeStock(w http.ResponseWriter, r *http.Request) {
	productName := r.URL.Query().Get("productName")
	if productName == "" {
		http.Error(w, "Missing productName parameter", http.StatusBadRequest)
		return
	}

	var request generated.PostChangeStockJSONBody

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "UPDATE products SET remaining_stock = $1 WHERE name = $2"
	_, err := dbConnection.Exec(query, request.NewStock, productName)
	if err != nil {
		http.Error(w, "Error updating stock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Stock changed successfully"})
}

func init() {
	http.HandleFunc("/lowStockProducts", StockHandler)
	http.HandleFunc("/incrementStock", IncrementStock)
	http.HandleFunc("/decrementStock", DecrementStock)
	http.HandleFunc("/changeStock", ChangeStock)
}
