package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	API_PATH = "/apis/v1/books"
)

type Library struct {
	dbHost, dbPass, dbName string
}

type Book struct {
	Id, Name, Isbn string
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "$haran@2468"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}

	l := Library{
		dbHost: dbHost,
		dbPass: dbPass,
		dbName: dbName,
	}

	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods("GET")
	r.HandleFunc(apiPath, l.postBook).Methods(http.MethodPost)
	http.ListenAndServe(":8080", r)
}

func (l *Library) getBooks(w http.ResponseWriter, r *http.Request) {
	log.Printf("getBooks called at ")

	db := l.openConnection()
	rows, err := db.Query("select * from books")
	if err != nil {
		fmt.Errorf("Error executing the Exec - %v", err)
	}
	var books []Book
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			fmt.Errorf("Error while scanning the row - %v", err)
		}
		b := Book{
			Id:   id,
			Name: name,
			Isbn: isbn,
		}
		books = append(books, b)
	}
	json.NewEncoder(w).Encode(books)
	l.closeConnection(db)

}

func (l *Library) postBook(w http.ResponseWriter, r *http.Request) {
	log.Printf("postBook called  ")
	db := l.openConnection()
	book := Book{}
	json.NewDecoder(r.Body).Decode(&book)
	//fmt.Println("Book after decoding is :", book)
	insertQuery, err := db.Prepare("insert into books values (?, ?, ?)")
	if err != nil {
		log.Fatalf("Error while preparing the statement - %v", err)
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error while starting the transaction - %v", err)
	}
	_, err = tx.Stmt(insertQuery).Exec(book.Id, book.Name, book.Isbn)
	if err != nil {
		fmt.Printf("Error while executing the statement - %v", err)
	}
	err = tx.Commit()
	if err != nil {
		fmt.Printf("Error while executing the statement - %v", err)
	}
	l.closeConnection(db)

}

func (l *Library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", l.dbPass, l.dbHost, l.dbName))
	if err != nil {
		panic(err)
	}
	return db
}

func (l *Library) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		panic(err)
	}
}
