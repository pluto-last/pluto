package pluto_bot

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	go http.ListenAndServe("localhost:9999", nil)

	// 断连后就重新连接
	for {
		Service()
		time.Sleep(time.Second * 10)
	}
}

func Service() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	err := Client(ctx)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}
}
