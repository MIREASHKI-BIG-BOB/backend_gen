package websocket

import (
	"backend_gen/internal/ports/websocket"
	gorillaWS "github.com/gorilla/websocket"
)

type client struct {
	conn *gorillaWS.Conn
	url  string
}

func (c *client) Connect(url string) error {
	conn, _, err := gorillaWS.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	c.url = url
	return nil
}

func (c *client) Disconnect() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		c.url = ""
		return err
	}
	return nil
}

func (c *client) IsConnected() bool {
	return c.conn != nil
}

func NewClient() websocket.Client {
	return &client{}
}
