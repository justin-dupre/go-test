package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func handleHello(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT * FROM public."2019"`)
	if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		http.Error(w, "Error fetching columns: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		for i, colName := range cols {
			fmt.Fprintf(w, "%s: %v\t", colName, columns[i])
		}
		fmt.Fprintln(w)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Row error: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	godotenv.Load()
	connStr := fmt.Sprintf("postgres://%s:%s@%s?sslmode=require", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_URL"))
	fmt.Printf("Connecting to database with connection string: %s\n", connStr)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	http.HandleFunc("/", handleHello)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
