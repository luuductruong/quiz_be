package grpc

import (
	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	coreGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/quiz_be/services/core/infra/config"
)

func NewClient(client *config.Client, initOptions ...coreGrpc.DialOption) (*coreGrpc.ClientConn, error) {
	retryOptions := []retry.CallOption{} // implement later
	defaultOption := []coreGrpc.DialOption{
		coreGrpc.WithChainUnaryInterceptor(retry.UnaryClientInterceptor(retryOptions...)),
	}
	if !client.SSL {
		defaultOption = append(defaultOption, coreGrpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	return coreGrpc.NewClient(client.Address(), append(initOptions, defaultOption...)...) // go 1.23 up
}
