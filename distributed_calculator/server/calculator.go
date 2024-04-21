package server

import (
	"context"
	"github.com/AkimKachaliev/Calculator_Golang---2-/server/grpc"
	"github.com/AkimKachaliev/Calculator_Golang---2-/server/http"
)

func Calculator(w http.ResponseWriter, r *http.Request) {
	// Implement the logic for handling calculator requests
}

func CalculatorGRPC(ctx context.Context, req *grpc.CalculatorRequest) (*grpc.CalculatorResponse, error) {
	// Implement the logic for handling calculator requests using gRPC
}
