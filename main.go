package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type server struct {
	subscriberMessageBuffer int
	subscribersMutex        sync.Mutex
	subscribers             map[*subscriber]struct{}
	mux                     http.ServeMux
}

type subscriber struct {
	msgs chan []byte
}

func NewServer() *server {
	s := &server{
		subscriberMessageBuffer: 10,
		subscribers:             make(map[*subscriber]struct{}, 0),
	}

	fs := http.FileServer(http.Dir("static/css"))
	s.mux.Handle("/", http.FileServer(http.Dir("htmx/")))
	s.mux.Handle("/css/", http.StripPrefix("/css/", fs))
	s.mux.HandleFunc("/ws", s.subscribeHandler)
	return s
}

func (s *server) subscribeHandler(writer http.ResponseWriter, req *http.Request) {
	err := s.subscribe(req.Context(), writer, req)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *server) addSubscriber(subscriber *subscriber) {
	s.subscribersMutex.Lock()
	s.subscribers[subscriber] = struct{}{}
	s.subscribersMutex.Unlock()
	fmt.Println("Added subscriber", subscriber)
}

func (s *server) subscribe(ctx context.Context, writer http.ResponseWriter, req *http.Request) error {
	var c *websocket.Conn
	subscriber := &subscriber{
		msgs: make(chan []byte, s.subscriberMessageBuffer),
	}
	s.addSubscriber(subscriber)

	c, err := websocket.Accept(writer, req, nil)
	if err != nil {
		return err
	}
	defer c.CloseNow()

	ctx = c.CloseRead(ctx)
	for {
		select {
		case msg := <-subscriber.msgs:
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			err := c.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				fmt.Println(err)
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func main() {
	fmt.Println("Hello, World!")

	server := NewServer()
	err := http.ListenAndServe(":8080", &server.mux)
	if err != nil {
		log.Fatal(err)
	}
}
