package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
}

func (app *application) googleAuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	uri := app.oauth.Google.Authorize("http://localhost:4000/auth/providers/google/callback")
	http.Redirect(w, r, uri, http.StatusSeeOther)
}

func (app *application) googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	m, err := app.oauth.Google.CompleteAuth(r, "http://localhost:4000/auth/providers/google/callback")
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	user, err := app.oauth.Google.FetchUser(m.AccessToken)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// register user

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(user)
}

func (app *application) errorResponse(w http.ResponseWriter, code int, err error) {
	log.Println(err)
	debug.Stack()
	http.Error(w, err.Error(), code)
}
