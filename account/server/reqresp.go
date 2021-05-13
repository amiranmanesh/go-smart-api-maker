package server

import (
	"context"
	"encoding/json"
	"net/http"
)

type (
	SignUpRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	SignUpResponse struct {
		Success bool   `json:"success"`
		Token   string `json:"token"`
	}
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	LoginResponse struct {
		Success bool   `json:"success"`
		Token   string `json:"token"`
	}
	VerifyRequest struct {
		Token string `json:"token"`
	}
	VerifyResponse struct {
		Success bool `json:"success"`
		UserID  uint `json:"user_id"`
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeSignUpReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeLoginReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeVerifyReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
