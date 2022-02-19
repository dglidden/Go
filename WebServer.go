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

func CreateUser(username string, password string, db *sql.DB) string {
	msg := "Creating user... "
	Query := "INSERT INTO USERS (username, password, created_at) values (?, ?, ?)"
	result, err := db.Exec(Query, username, password, time.Now())

	if err != nil {
		msg += fmt.Sprintf("There was an error creating your user: %s", err.Error())
	} else {
		UserID, err := result.LastInsertId()
		if err != nil {
			msg += fmt.Sprintf("There was an error retrieving your User ID: %s", err.Error())
		} else {
			msg += fmt.Sprintf("Your user was created with ID %d", UserID)
		}
	}

	return msg
}

func UpdatePassword(username string, password string, db *sql.DB) string {
	msg := "Updating password... "
	Query := "UPDATE USERS SET PASSWORD = ? WHERE USERNAME = ?"
	_, err := db.Exec(Query, password, username)
	if err != nil {
		msg += fmt.Sprintf("There was an error updating your password: %s", err.Error())
	} else {
		msg += "Updated password"
	}

	return msg
}

func main() {
	r := mux.NewRouter()

	db, err := sql.Open("mysql", "testuser:password@(fenris)/gotesting?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	r.HandleFunc("/db/{username}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		Query := "SELECT * FROM USERS WHERE username = ?"

		result, err := db.Query(Query, username)
		if err == nil {
			defer result.Close()
			if result.Next() {
				var id int64
				var password string
				var createdAt time.Time
				err := result.Scan(&id, &username, &password, &createdAt)
				if err != nil {
					fmt.Fprintf(w, "There was an error looking up your username %s", err.Error())
				} else {
					fmt.Fprintf(w, "The result of your query is: %d, %s, %s, %s", id, username, password, createdAt)
				}
			} else {
				fmt.Fprintf(w, "We were not able to find that username in the database")
			}
		} else {
			fmt.Fprintf(w, "There was an error looking up your username %s", err.Error())
		}
	})

	r.HandleFunc("/db/{username}/{password}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		password := vars["password"]

		Query := "SELECT * FROM USERS WHERE username = ?"

		result, err := db.Query(Query, username)
		if err == nil {
			if result.Next() {
				fmt.Fprint(w, UpdatePassword(username, password, db))
			} else {
				fmt.Fprint(w, CreateUser(username, password, db))
			}
		} else {
			fmt.Fprintf(w, "There was an error looking up your username %s", err.Error())
		}
	})

	http.ListenAndServe(":8080", r)
}
