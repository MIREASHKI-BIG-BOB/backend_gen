package websocket

import (
	"backend_gen/internal/ports/websocket"
	"fmt"
	"log/slog"
	"net/http"

	gorillaWS "github.com/gorilla/websocket"
)

type client struct {
	conn *gorillaWS.Conn
	url  string
}

func (c *client) Connect(url string, token string) error {
	slog.Info("Connecting to WebSocket server", "url", url, "token_length", len(token))

	headers := http.Header{}
	headers.Set("X-Auth-Sensor-Token", token)

	slog.Info("WebSocket headers prepared", "headers", headers)

	conn, resp, err := gorillaWS.DefaultDialer.Dial(url, headers)
	if err != nil {
		if resp != nil {
			slog.Error("Failed to dial WebSocket server",
				"url", url,
				"error", err,
				"status_code", resp.StatusCode,
				"status", resp.Status)
		} else {
			slog.Error("Failed to dial WebSocket server", "url", url, "error", err)
		}
		return err
	}

	c.conn = conn
	c.url = url
	slog.Info("WebSocket connection established", "url", url)
	return nil
}

func (c *client) Disconnect() error {
	if c.conn != nil {
		slog.Info("Disconnecting WebSocket", "url", c.url)
		err := c.conn.Close()
		c.conn = nil
		c.url = ""
		if err != nil {
			slog.Error("Error during WebSocket disconnect", "error", err)
		} else {
			slog.Info("WebSocket disconnected successfully")
		}
		return err
	}
	return nil
}

func (c *client) IsConnected() bool {
	return c.conn != nil
}

func (c *client) SendMessage(message []byte) error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	slog.Info("Sending WebSocket message", "message", string(message))
	err := c.conn.WriteMessage(gorillaWS.TextMessage, message)
	if err != nil {
		slog.Error("Failed to send WebSocket message", "error", err)
		return err
	}

	slog.Info("WebSocket message sent successfully")
	return nil
}

func NewClient() websocket.Client {
	return &client{}
}
