package main

import (
	"fmt"
	engineio "github.com/googollee/go-engine.io"
	"github.com/googollee/go-engine.io/transport"
	"github.com/googollee/go-engine.io/transport/polling"
	"github.com/googollee/go-engine.io/transport/websocket"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/cors"
	"log"
	"net/http"
)

var server *socketio.Server

func main() {
	pt := polling.Default
	wt := websocket.Default
	wt.CheckOrigin = func(req *http.Request) bool {
		return true
	}

	var err error
	server, err = socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{pt, wt},
	})
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "login", func(s socketio.Conn, msg string) {
		s.Join(msg)
		fmt.Println("login: ", msg)
	})
	server.OnError("/", func(e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})

	go server.Serve()
	defer server.Close()

	mux := http.NewServeMux()
	mux.Handle("/socket.io/", server)
	mux.HandleFunc("/p", publishMessage)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://dev.yichefu.cn", "https://dev.yichefu.cn"},
		AllowedMethods:   []string{"GET", "PUT", "OPTIONS", "POST", "DELETE"},
		AllowCredentials: true,
	})

	// decorate existing handler with cors functionality set in c
	handler := c.Handler(mux)

	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServeTLS(":8000", "fullchain.crt", "private.pem", handler))
}

func publishMessage(w http.ResponseWriter, r *http.Request) {
	//server.BroadcastToRoom(r.Form.Get("id"), "reply", "Broadcast")
	_ = r.ParseForm()
	server.BroadcastToRoom(r.Form.Get("id"), "message", "Broadcast", func(so socketio.Server, data string) {
		log.Println("Client ACK with data: ", data)
	})
	fmt.Fprintf(w, "BroadcastToRoom") //这个写入到w的是输出到客户端的
}
