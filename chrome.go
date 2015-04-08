package main

import (
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/websocket"
)

type chromeAdapter struct {
	server          *TriangleServer
	clientConnected bool
	rpcChan         chan int
}

func (a *chromeAdapter) name() string {
	return "Chrome"
}

func (a *chromeAdapter) load() {
	a.rpcChan = make(chan int)
	go a.serveWs()
}

func (a *chromeAdapter) toggle() {
	if a.clientConnected {
		a.rpcChan <- 0
	}
}

func (a *chromeAdapter) wsHandler(ws *websocket.Conn) {
	fmt.Println("Client connected")
	a.clientConnected = true
	// Read browser adapter state(is currently playing)
	var message string
	websocket.Message.Receive(ws, &message)
	currentlyPlaying, _ := strconv.ParseBool(string(message))
	fmt.Println("Is Playing", currentlyPlaying)
	if currentlyPlaying {
		a.server.setLastPlayingAdapter(a)
	}
	var readChan = make(chan string)
	go func() {
		for {
			var message string
			fmt.Println("Receiving")
			err := websocket.Message.Receive(ws, &message)
			if err != nil {
				ws.Close()
				readChan <- "close"
				break
			}
			readChan <- message
		}
	}()
	for {
		fmt.Println("selecting...")
		select {
		case <-a.rpcChan:
			fmt.Println("Got data from rpcChan")
			_, err := ws.Write([]byte("toggle"))
			if err != nil {
				fmt.Println("err", err)
				return
			}
			fmt.Println("Sent toggle")
		case message := <-readChan:
			switch message {
			case "close":
				fmt.Println("Closing")
				return
			default:
				fmt.Println("Got message", message)
				currentlyPlaying, _ := strconv.ParseBool(string(message))
				if currentlyPlaying {
					a.server.setLastPlayingAdapter(a)
				}
			}
		}
	}
}

func (a *chromeAdapter) serveWs() {
	http.Handle("/", websocket.Handler(a.wsHandler))
	fmt.Println("Serving WS at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServer: " + err.Error())
	}
}
