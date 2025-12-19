package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func UpdatePasswords(db *sql.DB) {
	rows, err := db.Query("SELECT id, password FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var oldPassword string
		if err := rows.Scan(&id, &oldPassword); err != nil {
			log.Fatal(err)
		}

		if len(oldPassword) < 50 {
			newHash, err := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Gagal generate hash user %d: %v\n", id, err)
				continue
			}

			_, err = db.Exec("UPDATE users SET password=$1 WHERE id=$2", string(newHash), id)
			if err != nil {
				log.Printf("Gagal update password user %d: %v\n", id, err)
				continue
			}

			fmt.Printf("Password user %d berhasil di-update\n", id)
		}
	}
}

func ResetAllPasswords(db *sql.DB, password string) error {
	rows, err := db.Query("SELECT id FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("Gagal scan user id: %v\n", err)
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Gagal generate hash user %d: %v\n", id, err)
			continue
		}

		_, err = db.Exec("UPDATE users SET password=$1 WHERE id=$2", string(hash), id)
		if err != nil {
			log.Printf("Gagal update password user %d: %v\n", id, err)
			continue
		}

		fmt.Printf("User %d: password berhasil di-reset\n", id)
	}

	return nil
}

func main() {
	db, err := sql.Open("postgres",
		"host=localhost port=5432 user=postgres password=1234567 dbname=authdb sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	defaultPassword := "12345"
	if err := ResetAllPasswords(db, defaultPassword); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Semua password user berhasil di-reset ke default:", defaultPassword)

	r := SetupRouter(db)

	fmt.Println("Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
