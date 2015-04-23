package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

/*
	Important: Enable Web-UI from within MPC-HC
	View > Options > Player > Web Interface
*/

const host = "localhost"
const port = 13579

type mpchcAdapter struct {
	server *TriangleServer
}

func (a *mpchcAdapter) name() string {
	return "MPC-HC"
}

func (a *mpchcAdapter) load() {
	mainUrl := fmt.Sprintf("http://%s:%d/", host, port)

	// Try to connect to web interface
	_, err := http.Get(mainUrl)
	if err != nil {
		fmt.Println("Error: couldn't connect to MPC-HC")
		return
	}

	// Query mpchc state and change Triangle lastPlayingAdapter if necessary
	state := getState()
	if state == "Playing" || a.server.lastPlayingAdapter == nil {
		a.server.setLastPlayingAdapter(a)
	}

	// TODO: state watcher (?)
}

func getState() string {
	varsUrl := fmt.Sprintf("http://%s:%d/variables.html", host, port)
	doc, err := goquery.NewDocument(varsUrl)
	if err != nil {
		fmt.Println("Error: couldn't connect to MPC-HC")
		return ""
	}

	return doc.Find("#statestring").Text()
}

// Toggles MPC-HC's playing state
func (a *mpchcAdapter) toggle() {
	doCmd("toggle")
}

// Plays previous track
func (a *mpchcAdapter) prev() {
	doCmd("prev")
}

// Plays next track
func (a *mpchcAdapter) next() {
	doCmd("next")
}

// Sends a command to MPC-HC's web interface
func doCmd(cmd string) {
	cmdUrl := fmt.Sprintf("http://%s:%d/command.html", host, port)
	cmdCodes := map[string]string{
		"toggle": "889",
		"prev":   "921",
		"next":   "922",
	}

	cmdCode := cmdCodes[cmd]

	_, err := http.PostForm(cmdUrl,
		url.Values{"wm_command": {cmdCode}, "submit": {"Go!"}})

	if err != nil {
		fmt.Println("Error: couldn't connect to MPC-HC")
	}
}
