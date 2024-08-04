package handler

import (
	"encoding/json"
	"github.com/NathanPr03/price-control/pkg/db"
	"net/http"
)

func DecrementStock(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}
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

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Stock decremented successfully"})
}

func init() {
	http.HandleFunc("/decrementStock", DecrementStock)
}
