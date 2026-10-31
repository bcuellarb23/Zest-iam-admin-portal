package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/bcuellarb23/Zest-iam-admin-portal/internal/routes"
)

func main() {
	godotenv.Load(".env")

	// FIXME: This handler will go to web, there I will have all web renders
	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			tmpl, err := template.ParseFiles("templates/dashboard.gohtml")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

	*/
	router := routes.SetupRouter()
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
