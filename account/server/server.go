package server

import (
	"context"
	"github.com/amiranmanesh/go-smart-api-maker/account/layers"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

//Transport Layer
func NewHTTPServer(ctx context.Context, endpoints layers.Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/signup").Handler(httptransport.NewServer(
		endpoints.SignUp,
		decodeSignUpReq,
		encodeResponse,
	))

	r.Methods("POST").Path("/login").Handler(httptransport.NewServer(
		endpoints.Login,
		decodeLoginReq,
		encodeResponse,
	))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
