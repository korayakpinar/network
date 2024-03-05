package ws

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/korayakpinar/p2pclient/src/price"
)

type Server struct {
	PriceCache *price.Cache
}

type GetPriceArgs struct {
	Pair price.Pair `json:"pair"`
}

type GetPriceReply struct {
	Price float64 `json:"price"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	for {
		var args GetPriceArgs

		err := ws.ReadJSON(&args)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		price := s.PriceCache.GetPrice(args.Pair)

		reply := GetPriceReply{Price: price}

		if err := ws.WriteJSON(reply); err != nil {
			log.Printf("error: %v", err)
			break
		}
	}
}

func (s *Server) Start(ctx context.Context, errChan chan error) {
	http.HandleFunc("/ws", s.handleConnections)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		errChan <- err
		return
	}
	log.Println("Websocket server started")
}

func NewServer(priceCache *price.Cache) *Server {
	return &Server{PriceCache: priceCache}
}
