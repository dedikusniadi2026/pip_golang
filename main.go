package main

import (
	"database/sql"
	"fmt"
	"log"
)

func main() {
	db, err := sql.Open("postgres",
		"host=localhost port=5432 user=postgres password=1234567 dbname=authdb sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	r := SetupRouter(db)

	fmt.Println("Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
