package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/bcuellarb23/Zest-iam-admin-portal/internal/auth"
	"github.com/bcuellarb23/Zest-iam-admin-portal/internal/store"
)

func main() {
	godotenv.Load(".env")

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

	http.HandleFunc("/signin", auth.SigninHandler())
	http.HandleFunc("/signout", auth.SignoutHandler)

	// FIXME: Call NewSessionStore to create a new empty session with the constructor, then Okta will return the values to populate that session.

	cur_session := store.NewSessionStore()
	_ = cur_session
	// Now Get data form okta

	/* token := auth.NewTokenGenerator()
	session := Session{
		UserName:  "GetUserFromOkta",
		Role:      "GetRoleFromOkta",
		ExpiresAt: time.Now().Add(time.Hour),
	} */

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
