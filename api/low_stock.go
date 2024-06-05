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
		_, _ = w.Write([]byte(fmt.Sprintf("Product %s is out of stock\n. Ordering more from central inventory service", product.ProductName)))
		time.Sleep(10 * time.Millisecond)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error with rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//things, _ := json.Marshal(lowStockProducts)

	err = pkg.SendEmail("Low stock products", fmt.Sprintf("The following products are low in stock: %v", string("some shit")))
	if err != nil {
		http.Error(w, "Error sending email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string][]Product{"products": lowStockProducts}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func init() {
	http.HandleFunc("/lowStockProducts", LowStock)
}
