package hazart

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// WSConn represents a simplified hijacked WebSocket connection
type WSConn struct {
	Conn net.Conn
	rw   *bufio.ReadWriter
}

// WriteMessage sends text message over socket
func (ws *WSConn) WriteMessage(msg []byte) error {
	_, err := ws.Conn.Write(msg)
	return err
}

// Close closes the underlying net.Conn
func (ws *WSConn) Close() error {
	return ws.Conn.Close()
}

// WebSocket upgrades standard HTTP connection to a WebSocket connection
func WebSocket(handler func(ws *WSConn)) HandlerFunc {
	return func(c *Context) {
		hj, ok := c.Writer.(http.Hijacker)
		if !ok {
			c.JSON(http.StatusInternalServerError, Map{"error": "webserver does not support hijacking"})
			return
		}

		conn, rw, err := hj.Hijack()
		if err != nil {
			c.JSON(http.StatusInternalServerError, Map{"error": fmt.Sprintf("hijack failed: %v", err)})
			return
		}

		// Write WebSocket handshake response
		rw.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
		rw.WriteString("Upgrade: websocket\r\n")
		rw.WriteString("Connection: Upgrade\r\n\r\n")
		rw.Flush()

		wsConn := &WSConn{Conn: conn, rw: rw}
		defer wsConn.Close()

		handler(wsConn)
	}
}
