package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type hero struct {
	Id       int
	Name     string
	Universe string
	Skill    string
	ImageUrl string
}

type villain struct {
	Id       int
	Name     string
	Universe string
	ImageUrl string
}

func listHeroes(db *sql.DB) ([]hero, error) {
	heroes := []hero{}

	query := `
	SELECT id, name, universe, skill, imageUrl from heroes
	`
	rows, err := db.Query(query)
	if err != nil {
		return []hero{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var h hero
		err = rows.Scan(&h.Id, &h.Name, &h.Universe, &h.Skill, &h.ImageUrl)
		if err != nil {
			return heroes, err
		}
		heroes = append(heroes, h)
	}

	// fmt.Println(heroes)
	return heroes, nil
}

func listVillains(db *sql.DB) ([]villain, error) {
	villains := []villain{}

	query := `
	SELECT id, name, universe, imageUrl from villain
	`
	rows, err := db.Query(query)
	if err != nil {
		return []villain{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var v villain
		err = rows.Scan(&v.Id, &v.Name, &v.Universe, &v.ImageUrl)
		if err != nil {
			return villains, err
		}
		villains = append(villains, v)
	}

	// fmt.Println(heroes)
	return villains, nil
}

func connect() (*sql.DB, error) {
	// Connect to the MySQL database
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/avenger")
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
	return db, nil
}

type product struct {
	Name string
	Type string
}

func main() {
	// connecting to db
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	heroes, err := listHeroes(db)
	fmt.Println(heroes)
	if err != nil {
		log.Fatal(err)
	}

	villains, err := listVillains(db)
	fmt.Println(villains)
	if err != nil {
		log.Fatal(err)
	}

	// create multiple endpoints using mux
	mux := http.NewServeMux()

	products := []product{
		{Name: "nasi goreng", Type: "makanan"},
		{Name: "nasi goreng", Type: "makanan"},
	}
	fmt.Println(products)

	mux.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	mux.HandleFunc("/heroes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(heroes)
	})

	mux.HandleFunc("/villains", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(villains)
		if err := json.NewEncoder(w).Encode(villains); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	// create web server
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	fmt.Println("Server running on port:8080")
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
