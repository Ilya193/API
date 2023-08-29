package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

type SaveReview struct {
	Title     string `json:"title"`
	ProductId int    `json:"productId"`
}

type Review struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	ProductId int    `json:"productId"`
}

func products(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@/products")

	result, err := db.Query("SELECT * FROM product")
	if err != nil {
		panic(err)
	}

	products := []Product{}

	defer db.Close()

	for result.Next() {
		p := Product{}
		result.Scan(&p.Id, &p.Name, &p.Description, &p.Url)
		products = append(products, p)
	}

	jsonResp, err := json.Marshal(products)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(jsonResp))
}

func saveReview(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@/products")

	if err != nil {
		panic(err)
	}

	var review SaveReview

	err = json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Decode error")
		return
	}

	_, err = db.Exec("insert into review (title, productId) values (?, ?)", review.Title, review.ProductId)
	if err != nil {
		panic(err)
	}

	result, err := db.Query("SELECT * FROM review where productId = ?", review.ProductId)
	if err != nil {
		panic(err)
	}

	reviews := []Review{}

	defer db.Close()

	for result.Next() {
		p := Review{}
		result.Scan(&p.Id, &p.Title, &p.ProductId)
		reviews = append(reviews, p)
	}

	jsonResp, err := json.Marshal(reviews)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(jsonResp))
}

func reviews(w http.ResponseWriter, r *http.Request) {
	productId := r.URL.Query().Get("productId")

	db, err := sql.Open("mysql", "root:root@/products")

	if err != nil {
		panic(err)
	}

	result, err := db.Query("SELECT * FROM review where productId = ?", productId)
	if err != nil {
		panic(err)
	}

	reviews := []Review{}

	defer db.Close()

	for result.Next() {
		p := Review{}
		result.Scan(&p.Id, &p.Title, &p.ProductId)
		reviews = append(reviews, p)
	}

	jsonResp, err := json.Marshal(reviews)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(jsonResp))
}

func main() {
	http.HandleFunc("/products", products)
	http.HandleFunc("/saveReview", saveReview)
	http.HandleFunc("/reviews", reviews)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
