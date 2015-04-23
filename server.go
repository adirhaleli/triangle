package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

// TriangleStatus contains information about the current state
// of Triangle
type TriangleStatus struct {
	LastPlayingAdapter string
}

// TriangleServer contains information about the server instance
type TriangleServer struct {
	adapters           []adapter
	lastPlayingAdapter adapter
}

// adapter interface defines the methods that Triangle adapters
// should implement
type adapter interface {
	name() string
	load()
	toggle()
}

// server function runs when the user starts the Triangle server
func server() {
	fmt.Println("Starting Triangle server...")
	serverInstance := &TriangleServer{}
	chrome := &chromeAdapter{server: serverInstance}
	mpd := &mpdAdapter{server: serverInstance}
	mpchc := &mpchcAdapter{server: serverInstance}
	serverInstance.adapters = []adapter{chrome, mpd, mpchc}

	// Load all adapters
	for _, adapter := range serverInstance.adapters {
		fmt.Printf("Loading adapter %s\n", adapter.name())
		adapter.load()
		fmt.Printf("Adapter %s loaded\n", adapter.name())
	}

	// Setup the RPC server for future Triangle requests
	serveRPC(serverInstance)
}

func serveRPC(serverInstance *TriangleServer) {
	rpc.Register(serverInstance)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Serving RPC at :1234")
	http.Serve(l, nil)
}

func (s *TriangleServer) setLastPlayingAdapter(adp adapter) {
	if s.lastPlayingAdapter == adp {
		return
	}
	s.lastPlayingAdapter = adp
	fmt.Println("Last playing adapter changed to", adp.name())
}

// --- RPC methods

// Toggle last playing adatper playing state
func (s *TriangleServer) Toggle(args *int, reply *int) error {
	if s.lastPlayingAdapter != nil {
		s.lastPlayingAdapter.toggle()
	} else {
		fmt.Println("lastPlayingAdapter is nil")
	}
	return nil
}

// Return info about the current state of Triangle
func (s *TriangleServer) Info(args *int, reply *TriangleStatus) error {
	// *reply = triangleStatus
	return nil
}
