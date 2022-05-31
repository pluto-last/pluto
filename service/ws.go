package service

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"pluto/model/params"
	"sync"
	"time"
)

var userClientMap UserClientMap

func init() {
	userClientMap = UserClientMap{Map: make(map[string]*UserClientPool)}
}

type WsService struct{}

func (wsService *WsService) GetUserClientMap() UserClientMap {
	return userClientMap
}

func (wsService *WsService) GetUserClientPool() *UserClientPool {
	return new(UserClientPool)
}

// Deprecated: Use WSPack instead.
type WSResp struct {
	Type string
	Code int
	Data json.RawMessage
	Msg  string
}

type WSPack struct {
	TaskID string
	Type   string
	Code   int
	Data   interface{}
	Msg    string
}

func WSPackOK(Type string, data interface{}) (*WSPack, error) {
	return &WSPack{Type: Type, Data: data}, nil
}

type WsConn struct {
	mutex sync.Mutex
	*websocket.Conn
}

type UserClientMap struct {
	Map   map[string]*UserClientPool
	mutex sync.Mutex
}

type UserClientPool struct {
	ID            string
	DeleteID      string // 删除时核对对象
	Ctx           context.Context
	CancelFunc    context.CancelFunc
	Conn          *WsConn
	HeartbeatTime int64 // 最后心跳时间,单位秒
}

type Auth struct {
	Time time.Time
	Md5  string
}

func (userClientMap UserClientMap) Set(poolID string, clientPool *UserClientPool) {
	userClientMap.mutex.Lock()
	defer userClientMap.mutex.Unlock()
	userClientMap.Map[poolID] = clientPool
}

func (userClientMap UserClientMap) Get(poolID string) (*UserClientPool, bool) {
	userClientMap.mutex.Lock()
	defer userClientMap.mutex.Unlock()
	ctl, ok := userClientMap.Map[poolID]
	return ctl, ok
}

func (userClientMap UserClientMap) Delete(poolID string, deleteID string) {
	userClientMap.mutex.Lock()
	defer userClientMap.mutex.Unlock()
	pool, ok := userClientMap.Map[poolID]
	if ok && pool.DeleteID == deleteID {
		delete(userClientMap.Map, poolID)
	}
}

// UserClientHeartbeat 心跳检测
func (clientPool *UserClientPool) UserClientHeartbeat(ctx context.Context) {
	t := time.NewTicker(time.Minute)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-clientPool.Ctx.Done():
			return
		case <-t.C:
			// 心跳测试
			err := clientPool.Conn.WriteJSON(params.WSPack{
				Type: "Heart",
				Code: 0,
				Msg:  "heartbeat",
			})

			if err != nil {
				clientPool.CancelFunc()
				return
			}
		}
	}
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

func (clientPool *UserClientPool) Close() {
	clientPool.CancelFunc()
	clientPool.Conn.Close()
	userClientMap.Delete(clientPool.ID, clientPool.DeleteID)
}

func (clientPool *UserClientPool) Handler(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			clientPool.Conn.Close()
			return nil
		case <-clientPool.Ctx.Done():
			clientPool.Conn.Close()
			return nil
		default:
		}

		t, data, err := clientPool.Conn.ReadMessage()
		if err != nil {
			clientPool.CancelFunc()
			return err
		}

		switch t {
		case websocket.PingMessage:
		case websocket.PongMessage:
		case websocket.TextMessage, websocket.BinaryMessage:
			pack := new(WSPack)
			json.Unmarshal(data, pack)

			if pack.Type == "ping" {
				resp := WSPack{Type: "pong"}
				clientPool.HeartbeatTime = time.Now().Unix()
				err = clientPool.Conn.WriteJSON(&resp)
			}
		}
	}
}
