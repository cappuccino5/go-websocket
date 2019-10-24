package chat

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type Conn struct {
	wConn     *websocket.Conn
	closeChan chan struct{}
	pChan     chan struct{}

	// Client should ping server in right duration
	cTimer        *time.Timer
	receiveHandle func(data []byte, msgType int)
	closeHandle   func()
}

const (
	// content of ping
	_PING_CONTENT = "[ping]:status"
	// timeout of ping
	_PING_TIMEOUT = 5 * time.Second
	// content of heartbeats
	_HEARTBEAT_CONTENT = "[pong]:heartbeat"
	// timeout of heartbeats
	_HB_TIMEOUT = 60 * 5 * time.Second
)

func newConn(wsConn *websocket.Conn) *Conn {
	conn := &Conn{
		wConn:     wsConn,
		pChan:     make(chan struct{}),
		closeChan: make(chan struct{}),
		cTimer:    time.NewTimer(_HB_TIMEOUT),
	}

	conn.wConn.SetCloseHandler(func(code int, text string) error {
		conn.close()
		conn.wConn.WriteControl(websocket.CloseMessage, []byte(""), time.Now().Add(time.Second))
		return nil
	})

	conn.wConn.SetPongHandler(func(appData string) error {
		if appData == _PING_CONTENT {
			conn.pChan <- struct{}{}
		}
		return nil
	})

	conn.wConn.SetPingHandler(func(appData string) error {
		err := conn.wConn.WriteControl(websocket.PongMessage, []byte(_HEARTBEAT_CONTENT), time.Now().Add(time.Second))
		conn.cTimer.Reset(_HB_TIMEOUT)
		return err
	})

	go conn.readLoop()

	// listen to Client's heart beat
	go func() {
		<-conn.cTimer.C
		fmt.Println("heartbeats timeout")
		conn.close()
	}()

	return conn
}

func (c *Conn) sendMessage(messageType int, data []byte) error {
	go c.ping()
	timer := time.NewTimer(_PING_TIMEOUT)
	select {
	case <-timer.C:
		return ERR01
	case <-c.pChan:
		timer.Stop()
	}
	if err := c.wConn.WriteMessage(messageType, data); err != nil {
		return err
	}
	return nil
}

func (c *Conn) ping() {
	c.wConn.WriteControl(websocket.PingMessage, []byte(_PING_CONTENT), time.Now().Add(time.Second))
}

func (c *Conn) readLoop() {
ReadLoop:
	for {
		t, p, err := c.wConn.ReadMessage()
		if err != nil {
			break ReadLoop
		}

		fmt.Printf("Receve Data: %s, \n", p)

		select {
		case <-c.closeChan:
			break ReadLoop
		default:
			c.receiveHandle(p, t)
		}
	}
}

func (c *Conn) close() {
	select {
	case <-c.closeChan:
	default:
		c.wConn.Close()
		c.closeHandle()
		c.cTimer.Stop()
		close(c.closeChan)
	}
}
