package client

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/loamhoof/indicator"
)

type ShepherdClient struct {
	port   int
	conn   *grpc.ClientConn
	client pb.ShepherdClient
}

func NewShepherdClient(port int) *ShepherdClient {
	return &ShepherdClient{port: port}
}

func (sc *ShepherdClient) Init() error {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", sc.port), grpc.WithInsecure())
	if err != nil {
		return err
	}

	sc.conn = conn
	sc.client = pb.NewShepherdClient(conn)

	return nil
}

func (sc *ShepherdClient) Close() {
	sc.conn.Close()
}

func (sc *ShepherdClient) Update(req *pb.Request) (*pb.Response, error) {
	return sc.client.Update(context.Background(), req)
}
