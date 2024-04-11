package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PriyanshuSharma23/oauth-demo/internals/oauth"
)

type application struct {
	oauth *oauth.Oauth
}

func main() {
	mux := http.NewServeMux()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	oa, err := oauth.New()
	if err != nil {
		log.Fatal("failed to load env variables")
	}

	app := &application{
		oauth: oa,
	}

	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/auth/providers/google/authorize", app.googleAuthorizeHandler)
	mux.HandleFunc("/auth/providers/google/callback", app.googleCallbackHandler)

	log.Println("Starting server at port 4000")
	http.ListenAndServe(fmt.Sprintf(":%d", 4000), mux)
}
