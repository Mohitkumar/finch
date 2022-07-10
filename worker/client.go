package worker

import (
	"google.golang.org/grpc"

	api "github.com/mohitkumar/finch/api/v1"
)

type client struct {
	conn *grpc.ClientConn
}

func NewClient(serverAddress string) (*client, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &client{
		conn: conn,
	}, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) GetApiClient() api.TaskServiceClient {
	return api.NewTaskServiceClient(c.conn)
}
