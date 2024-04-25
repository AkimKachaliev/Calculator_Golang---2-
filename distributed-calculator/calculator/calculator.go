package calculator

import (
	"context"
	"github.com/AkimKachaliev/distributed-calculator/models"
	"google.golang.org/grpc"
	"net"
)

type CalculatorServer struct {
	models.UnimplementedCalculatorServer
}

func (s *CalculatorServer) Calculate(ctx context.Context, req *CalculateRequest) (*CalculateResponse, error) {
	expression := models.Expression{
		UserID:     req.GetUserID(),
		Expression: req.GetExpression(),
	}

	err := expression.Calculate()
	if err != nil {
		return nil, err
	}

	return &CalculateResponse{Result: expression.Result}, nil
}

func RunGRPCServer() error {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	models.RegisterCalculatorServer(s, &CalculatorServer{})

	return s.Serve(lis)
}
