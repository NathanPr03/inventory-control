package handler

import (
	"encoding/json"
	"github.com/NathanPr03/price-control/pkg/db"
	"net/http"
)

func StockHandler(w http.ResponseWriter, r *http.Request) {
	dbConnection, err := db.ConnectToDb()
	if err != nil {
		http.Error(w, "Error connecting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

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

func init() {
	http.HandleFunc("/lowStockProducts", StockHandler)
	http.HandleFunc("/incrementStock", IncrementStock)
	http.HandleFunc("/decrementStock", DecrementStock)
	http.HandleFunc("/changeStock", ChangeStock)
}
