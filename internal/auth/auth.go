package auth

import (
	"net/http"
	"os"

	oktaUtils "github.com/okta/samples-golang/okta-hosted-login/utils"
)

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "no-cache")

	nonce, _ := oktaUtils.GenerateNonce()
	state, _ := GenerateState()

	var redirectPath string

	q := r.URL.Query()
	q.Add("client_id", os.Getenv("OKTA_CLIENT_ID"))
	q.Add("response_type", "code")
	q.Add("response_mode", "query")
	q.Add("scope", "openid profile email")
	q.Add("redirect_uri", os.Getenv("BASE_URL")+"/authorization-code/callback")
	q.Add("state", state)
	q.Add("nonce", nonce)

	redirectPath = os.Getenv("ISSUER") + "v1/authorize?" + q.Encode()

	http.Redirect(w, r, redirectPath, http.StatusFound)
}

func SignoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "Okta-hosted-signin-session-store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	delete(session.Values, "id_token")
	delete(session.Values, "access_token")

	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
