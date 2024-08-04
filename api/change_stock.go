package handler

import (
	"encoding/json"
	"github.com/NathanPr03/price-control/pkg/db"
	"inventory-control/api/generated"
	"net/http"
)

func ChangeStock(w http.ResponseWriter, r *http.Request) {
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

	var request generated.PostChangeStockJSONBody

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "UPDATE products SET remaining_stock = $1 WHERE name = $2"
	_, err = dbConnection.Exec(query, request.NewStock, productName)
	if err != nil {
		http.Error(w, "Error updating stock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Stock changed successfully"})
}

func init() {
	http.HandleFunc("/changeStock", ChangeStock)
}
