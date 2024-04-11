package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PriyanshuSharma23/oauth-demo/internals/oauth"
	oauthgoth "github.com/PriyanshuSharma23/oauth-demo/internals/oauth-goth"
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

type application struct {
	oauth *oauth.Oauth
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	oa, err := oauth.New()
	if err != nil {
		log.Fatal("failed to load env variables")
	}

	oauthgoth.InitializeGoth()

	app := &application{
		oauth: oa,
	}

	router := chi.NewRouter()

	router.Get("/", app.homeHandler)

	router.Get("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
			w.Header().Set("Content-Type", "application/json")
			enc := json.NewEncoder(w)
			enc.SetIndent("", "\t")
			enc.Encode(gothUser)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	})

	router.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "\t")
		enc.Encode(user)
	})

	router.Get("/logout/{provider}", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

		gothic.Logout(w, r)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	log.Println("Starting server at port 4000")
	http.ListenAndServe(fmt.Sprintf(":%d", 4000), router)
}
