package main

type ErrorResponse struct {
	Reason string `json:"reason"`
	Code   int    `json:"code"`
}
