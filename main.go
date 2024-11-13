package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/connect", connectHandler)

	port := ":80"
	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Someone hit me!")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h2>Wise Words</h2><p>Keep ya head up</p>")
}

func getDBConnection() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "database"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPass == "" {
		dbPass = "postgres"
	}
	if dbName == "" {
		dbName = "postgres"
	}

	connStr := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName,
	)
	fmt.Println(connStr)

	var db *sql.DB
	var err error
	for i := 0; i < 3; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Attempt %d: Failed to open database: %v", i+1, err)
			time.Sleep(time.Second * 2)
			continue
		}

		err = db.Ping()
		if err == nil {
			return db, nil
		}
		log.Printf("Attempt %d: Failed to ping database: %v", i+1, err)
		time.Sleep(time.Second * 2)
	}
	return nil, fmt.Errorf("failed to connect to database after 3 attempts: %v", err)
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	db, err := getDBConnection()
	if err != nil {
		log.Printf("Failed connect to db: %v", err)
		fmt.Fprintf(w, "<h1>Failed to connect to database</h1><p>Error: %v</p>", err)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO dimassetra_access_log (timestamp) VALUES ($1)", time.Now())
	if err != nil {
		log.Printf("Failed to insert data: %v", err)
		fmt.Fprintf(w, "<h1>Failed to insert data</h1><p>Error: %v</p>", err)
		return
	}

	log.Println("Success connect to db")
	fmt.Fprintf(w, "<h1>Successfully connected to database</h1>")
}
