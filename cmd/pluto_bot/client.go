package pluto_bot

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"golang.org/x/xerrors"
	"io/ioutil"
	"log"
	"net/http"
	"pluto/global"
	"pluto/utils"
	"sync"
	"time"
)

type WsConn struct {
	mutex sync.Mutex
	*websocket.Conn
}

func (c *WsConn) WriteJSON(v interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Conn.WriteJSON(v)
}

func (c *WsConn) WriteMessage(messageType int, data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Conn.WriteMessage(messageType, data)
}

func Client(ctx context.Context) error {
	log.Println("Service")

	subctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	path := global.GVA_CONFIG.WsPath

	// 发起webscoket链接
	wsconn, resp, err := websocket.DefaultDialer.DialContext(subctx, path, http.Header{})
	if err != nil {
		if resp == nil || resp.StatusCode != 200 {
			return xerrors.Errorf("%w", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return xerrors.Errorf("%w", err)
		}
		err = xerrors.New(string(body))
		if err != nil {
			log.Printf("%+v\n", err)
		}
		return err
	}
	defer wsconn.Close()

	// webscoket 已经连接
	log.Println("connected")
	conn := &WsConn{
		Conn: wsconn,
	}

	// 心跳检测
	go func() {
		for {
			select {
			case <-subctx.Done():
				return
			default:
			}
			time.Sleep(time.Second * 30)
			log.Println("ping")
			err = conn.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				log.Printf("%+v\n", err)
				return
			}
		}
	}()

	// 通信处理
	for {
		select {
		case <-subctx.Done():
			return nil
		default:
		}
		t, data, err := conn.ReadMessage()
		if err != nil {
			return xerrors.Errorf("%w", err)
		}
		switch t {
		case websocket.PingMessage, websocket.PongMessage:
		case websocket.TextMessage, websocket.BinaryMessage:
			go Handler(subctx, conn, data)
		}
	}
}

// Handler 处理具体的通信数据
func Handler(ctx context.Context, conn *WsConn, data []byte) {
	log.Println("Got Websocket Data:", string(data))
	pack := new(utils.WSPack)
	err := json.Unmarshal(data, pack)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}

}
