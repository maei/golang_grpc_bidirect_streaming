package server

import (
	"fmt"
	"github.com/maei/golang_grpc_bidirect_streaming/grpc_server/src/domain/greetpb"
	"github.com/maei/shared_utils_go/logger"
	"google.golang.org/grpc"
	"io"
	"net"
)

type server struct{}

var (
	s = grpc.NewServer()
)

func (*server) GetGreeting(stream greetpb.GreetService_GetGreetingServer) error {
	logger.Info("gRPC greet-streaming started")

	for {
		req, reqErr := stream.Recv()
		if reqErr == io.EOF {
			return nil
		}
		if reqErr != nil {
			logger.Error("Error while fetich GRPC-Client request", reqErr)
			return reqErr
		}
		res := &greetpb.GreetResponse{
			Result: fmt.Sprintf("Hello %v from the GRPC-Server", req.GetGreet().GetFirstName()),
		}
		streamErr := stream.Send(res)
		if streamErr != nil {
			logger.Error("Error while streaming data to GRPC-Client", streamErr)
			return streamErr
		}
	}
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("error while listening gRPC Server", err)
	}

	greetpb.RegisterGreetServiceServer(s, &server{})

	errServer := s.Serve(lis)
	if errServer != nil {
		logger.Error("error while serve gRPC Server", errServer)
	}
}
