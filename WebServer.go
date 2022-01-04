package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	db, err := sql.Open("mysql", "testuser:password@(fenris)/gotesting?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	r.HandleFunc("/db/{username}/{paassword}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		password := vars["password"]

		Query := `INSERT INTO users (username, password, created_at) values (?, ?, ?)`

		result, err := db.Exec(Query, username, password, time.Now())
		if err != nil {
			fmt.Fprintf(w, "There was an error creating your user: %s", err.Error())
		} else {
			UserID, err := result.LastInsertId()
			if err != nil {
				fmt.Fprintf(w, "There was an error retrieving your User ID: %s", err.Error())
			} else {
				fmt.Fprintf(w, "Your user was created with ID %d", UserID)
			}
		}
	})

	http.ListenAndServe(":8080", r)
}
