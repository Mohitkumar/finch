package worker

import (
	"google.golang.org/grpc"

	api "github.com/mohitkumar/finch/api/v1"
)

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(serverAddress string) (*Client, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetApiClient() api.TaskServiceClient {
	return api.NewTaskServiceClient(c.conn)
}
