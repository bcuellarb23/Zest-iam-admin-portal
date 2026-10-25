package auth

import (
	"net/http"
	"os"

	"github.com/bcuellarb23/Zest-iam-admin-portal/internal/store"
	oktaUtils "github.com/okta/samples-golang/okta-hosted-login/utils"
)

// FIXME: Use Okta Developer web page to fix this handler, Remember we are using OIDC for authentication. remember also grabbing data from the .env
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
