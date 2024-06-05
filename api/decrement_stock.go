package handler

import (
	"encoding/json"
	"github.com/NathanPr03/price-control/pkg/db"
	"net/http"
)

func DecrementStock(w http.ResponseWriter, r *http.Request) {
	dbConnection, err := db.ConnectToDb()
	if err != nil {
		http.Error(w, "Error connecting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	productName := r.URL.Query().Get("productName")
	if productName == "" {
		http.Error(w, "Missing productName parameter", http.StatusBadRequest)
		return
	}

	query := "UPDATE products SET remaining_stock = remaining_stock - 1 WHERE name = $1"
	_, err = dbConnection.Exec(query, productName)
	if err != nil {
		http.Error(w, "Error updating stock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Stock decremented successfully"})
}

func init() {
	http.HandleFunc("/decrementStock", DecrementStock)
}
