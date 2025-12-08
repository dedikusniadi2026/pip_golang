package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestRouteExists(t *testing.T) {
	db, _ := sql.Open("postgres", "postgres://postgres:1234567@localhost:5432/authdb?sslmode=disable")

	r := SetupRouter(db)

	req, _ := http.NewRequest("GET", "/drivers", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, 404, w.Code)
}
