package server

import (
	"context"
	"github.com/amiranmanesh/go-smart-api-maker/account/proto"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

//Transport Layer
func NewHTTPServer(endpoints Endpoints) http.Handler {
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

type grpcServer struct {
	verifyTokenHandler grpctransport.Handler
}

func NewGRPCServer(ctx context.Context, endpoints Endpoints) proto.AccountServiceServer {
	return &grpcServer{
		verifyTokenHandler: grpctransport.NewServer(
			endpoints.VerifyToken,
			decodeVerifyTokenRequest,
			encodeVerifyTokenRequest,
		),
	}
}

func (s *grpcServer) VerifyToken(ctx context.Context, req *proto.VerifyTokenRequest) (*proto.VerifyTokenResponse, error) {
	_, response, err := s.verifyTokenHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return response.(*proto.VerifyTokenResponse), nil
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
