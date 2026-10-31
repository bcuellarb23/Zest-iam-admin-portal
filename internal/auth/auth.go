package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/bcuellarb23/Zest-iam-admin-portal/internal/store"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
)

// Define the environment as recommended by okta, Didn't use client secret since I am using a pair public/private key
// This config will be in a function, It didn't worked out of a function because it was configurtes before main() loaded .env
func oauthConfig() (*oauth2.Config, error) {
	clientID := strings.TrimSpace(os.Getenv("OKTA_CLIENT_ID"))
	issuer := strings.TrimRight(strings.TrimSpace(os.Getenv("OKTA_ISSUER")), "/")
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("BASE_URL")), "/")
	if clientID == "" || issuer == "" || baseURL == "" {
		return nil, errors.New("missing required Okta OAuth configuration")
	}

	return &oauth2.Config{
		ClientID:    clientID,
		RedirectURL: baseURL + "/authorization-code/callback",
		Scopes:      []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  issuer + "/v1/authorize",
			TokenURL: issuer + "/v1/token",
		},
	}, nil
}

// Use Okta Developer web page for this handler, Remember we are using OIDC for authentication. remember also grabbing data from the .env
func SigninHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Use PKCE to protect against CSRF attacks
	verifier := oauth2.GenerateVerifier()
	state, _ := GenerateState()

	// Store state in a short lived http cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	// Store the verifier in a short lived http cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_verifier",
		Value:    verifier,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	// Have to call the function to gen the config
	conf, err := oauthConfig()
	if err != nil {
		http.Error(w, "Okta OAuth configuration is missing OKTA_CLIENT_ID, OKTA_ISSUER, or BASE_URL", http.StatusInternalServerError)
		return
	}
	// Generate the authorization URL including the PKCE challenge code
	authURL := conf.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(verifier),
	)

	// Redirect the user's browser to Okta
	http.Redirect(w, r, authURL, http.StatusFound)
}

func SignoutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cur_session := store.NewSessionStore()
	session, err := cur_session.Get("Okta-hosted-signin-session-store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_ = session
	/* delete(//session.Values, "id_token")
	delete(//session.Values, "access_token")

	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound) */
}
