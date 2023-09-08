package authsvc

import (
	"fmt"
	"gateway/config"
	"gateway/pb"

	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.UserControllerClient
}

func InitServiceClient(c *config.Config) pb.UserControllerClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.AuthSvcUrl, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewUserControllerClient(cc)
}
