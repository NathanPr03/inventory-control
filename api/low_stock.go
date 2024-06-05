package handler

import (
	"encoding/json"
	"github.com/NathanPr03/price-control/pkg/db"
	"net/http"
)

type Product struct {
	ProductName    string `json:"productName"`
	RemainingStock int    `json:"remainingStock"`
}

func LowStock(w http.ResponseWriter, r *http.Request) {
	dbConnection, err := db.ConnectToDb()
	if err != nil {
		http.Error(w, "Error connecting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	query := "SELECT name, remaining_stock FROM products WHERE remaining_stock < 10"
	rows, err := dbConnection.Query(query)
	if err != nil {
		http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ProductName, &product.RemainingStock); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error with rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string][]Product{"products": products}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func init() {
	http.HandleFunc("/lowStockProducts", LowStock)
}
