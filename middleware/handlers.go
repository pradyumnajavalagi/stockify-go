package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"psq-project/models"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Fatal("Error opening database connection", err)
	}

	// Test connection
	if err = db.Ping(); err != nil {
		log.Fatal("Database connection failed", err)
	}

	fmt.Println("Successfully connected to PostgreSQL")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var stock models.Stock
	if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	insertID := insertStock(stock)
	res := response{
		ID:      insertID,
		Message: "Stock created successfully",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	stock, err := getStock(int64(id))
	if err != nil {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(stock)
}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stocks, err := getStocks()
	if err != nil {
		http.Error(w, "Failed to fetch stocks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	var stock models.Stock
	if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updatedRows, err := updateStock(int64(id), stock)
	if err != nil {
		http.Error(w, "Failed to update stock", http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Update successful: %v rows affected", updatedRows)
	res := response{ID: int64(id), Message: msg}
	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	deletedRows := deleteStock(int64(id))
	msg := fmt.Sprintf("Deletion successful: %v rows affected", deletedRows)
	res := response{ID: int64(id), Message: msg}
	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $3) RETURNING stocksid`
	var id int64
	if err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted stock with ID %v\n", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `SELECT stocksid, name, price, company FROM stocks WHERE stockid = $1`
	var stock models.Stock
	row := db.QueryRow(sqlStatement, id)
	if err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company); err != nil {
		if err == sql.ErrNoRows {
			return stock, fmt.Errorf("Stock not found")
		}
		return stock, err
	}
	return stock, nil
}

func getStocks() ([]models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `SELECT stocksid, name, price, company FROM stocks`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func updateStock(id int64, stock models.Stock) (int64, error) {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `UPDATE stocks SET name = $2, price = $3, company = $4 WHERE stocksid = $1`
	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	fmt.Printf("%v rows affected\n", rowsAffected)
	return rowsAffected, nil
}

func deleteStock(id int64) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stocks WHERE stocksid = $1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v rows deleted\n", rowsAffected)
	return rowsAffected
}
