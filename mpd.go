package main

import (
	"fmt"

	"github.com/fhs/gompd/mpd"
)

type mpdAdapter struct {
	server *TriangleServer
	client *mpd.Client
}

func (a *mpdAdapter) name() string {
	return "MPD"
}

func (a *mpdAdapter) load() {
	// Connect to MPD
	var err error
	a.client, err = mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		fmt.Println("Error: couldn't connect to MPD")
		return
	}

	// Query MPD state and change Triangle lastPlayingAdapter if necessary
	status, _ := a.client.Status()
	if status["state"] == "play" || a.server.lastPlayingAdapter == nil {
		a.server.setLastPlayingAdapter(a)
	}

	// Setup watcher to be notified when MPD state change
	watcher, err := mpd.NewWatcher("tcp", ":6600", "")
	if err != nil {
		fmt.Println("Error: couldn't connect to MPD")
		return
	}
	go func() {
		for subsystem := range watcher.Event {
			fmt.Println("Changed subsystem", subsystem)
			if subsystem == "player" {
				status, _ := a.client.Status()
				if status["state"] == "play" {
					a.server.setLastPlayingAdapter(a)
				}
			}
		}
	}()
}

// Toggles MPD playing state
func (a *mpdAdapter) toggle() {
	status, _ := a.client.Status()
	shouldPause := status["state"] == "play"
	a.client.Pause(shouldPause)
}

// func queryMPD() {
// 	conn, err := mpd.Dial("tcp", "localhost:6600")
// 	if err != nil {
// 		fmt.Println("Error: couldn't connect to MPD")
// 		return
// 	}
// 	defer conn.Close()
// 	status, _ := conn.Status()
// 	if status["state"] == "play" {
// 		// triangleStatus.LastPlayingDevice = "mpd"
// 	}
// }
