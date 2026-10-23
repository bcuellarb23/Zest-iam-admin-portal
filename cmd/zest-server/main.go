package main

import (
	"net/http"
	"time"

	oktaUtils "github.com/okta/samples-golang/okta-hosted-login/utils"

	"github.com/bcuellarb23/Zest-iam-admin-portal/internal/auth"
	"github.com/bcuellarb23/Zest-iam-admin-portal/internal/store"
)

func main() {
	oktaUtils.ParseEnvironment()

	http.HandleFunc("/signin", auth.SigninHandler)
	http.HandleFunc("/signout", auth.SignoutHandler)

	userStore := store.NewSessionStore()

	token := auth.NewTokenGenerator()
	session := Session{
		UserName:  "GetUserFromOkta",
		Role:      "GetRoleFromOkta",
		ExpiresAt: time.Now().Add(time.Hour),
	}
}
