package calculator

import (
	"context"
	"net"

	"github.com/AkimKachaliev/distributed-calculator/distributed-calculator/models"
	"google.golang.org/grpc"
	"github.com/AkimKachaliev/distributed-calculator/distributed-calculator/api/v1/calculator"
)

type CalculatorServer struct {
	calculator.UnimplementedCalculatorServer
}

func (s *CalculatorServer) Calculate(ctx context.Context, req *calculator.CalculateRequest) (*calculator.CalculateResponse, error) {
	expression := models.Expression{
		UserID:     req.GetUserID(),
		Expression: req.GetExpression(),
	}

	err := expression.Calculate()
	if err != nil {
		return nil, err
	}

	return &calculator.CalculateResponse{Result: expression.Result}, nil
}

func RunGRPCServer() error {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	calculator.RegisterCalculatorServer(s, &CalculatorServer{})

	return s.Serve(lis)
}
