package server

import (
	"context"
	"github.com/amiranmanesh/go-smart-api-maker/account/proto"
	gt "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

type gRPCServer struct {
	verifyTokenHandler gt.Handler
}

//Transport Layer
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
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

func NewGRPCServer(ctx context.Context, endpoints Endpoints) proto.AccountServiceServer {

	return &gRPCServer{
		verifyTokenHandler: gt.NewServer(
			endpoints.VerifyToken,
			decodeVerifyTokenRequest,
			encodeVerifyTokenRequest,
		),
	}
}

func (s *gRPCServer) VerifyToken(ctx context.Context, req *proto.VerifyTokenRequest) (*proto.VerifyTokenResponse, error) {
	_, response, error := s.verifyTokenHandler.ServeGRPC(ctx, req)
	if error != nil {
		return nil, error
	}

	return response.(*proto.VerifyTokenResponse), nil
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
