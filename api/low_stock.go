package handler

import (
	"encoding/json"
	"fmt"
	"github.com/NathanPr03/price-control/pkg/db"
	"inventory-control/pkg"
	"net/http"
	"time"
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

	var noStockProducts []Product
	var lowStockProducts []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ProductName, &product.RemainingStock); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		lowStockProducts = append(lowStockProducts, product)

		if product.RemainingStock < 1 {
			noStockProducts = append(noStockProducts, product)
		}
	}

	for _, product := range noStockProducts {
		println(fmt.Sprintf("Product %s is out of stock. Ordering more from central inventory service", product.ProductName))
		time.Sleep(10 * time.Millisecond)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error with rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = pkg.SendEmail("Low stock products", prettyPrintProducts(lowStockProducts))
	if err != nil {
		http.Error(w, "Error sending email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string][]Product{"products": lowStockProducts}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func prettyPrintProducts(products []Product) string {
	var result string

	result += "Low Stock Products:\n"
	result += "-------------------------------\n"
	result += fmt.Sprintf("%-20s | %-15s\n", "Product Name", "Remaining Stock")
	result += "-------------------------------\n"
	for _, product := range products {
		result += fmt.Sprintf("%-20s | %-15d\n", product.ProductName, product.RemainingStock)
	}
	result += "-------------------------------\n"

	return result
}

func init() {
	http.HandleFunc("/lowStockProducts", LowStock)
}
