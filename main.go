package main

import (
	"fmt"
	"os"
	"github.com/zserge/webview"
	"github.com/pebbe/zmq4"
	"encoding/json"
	"github.com/y0ssar1an/q"
)

func parseCommand(msg []byte) Command {
	var v []interface{}
	var m Command
	json.Unmarshal(msg, &v)
	
	kind := v[0].(string)
	if kind == "Sub" {
		m = NewSubCommand(v)
	} else if kind == "PatchAttrs" {
		m = NewPatchAttrsCommand(v)
	} else if kind == "PostHtml" {
		m = NewPostHTMLCommand(v)
	} else if kind == "PatchCss" {
		m = NewPatchCSSCommand(v)
	} else if kind == "Close" {
		m = CloseCommand{}
	} else {
		q.Q(v)
	}
	return m
}



type Guise struct {
	View webview.WebView
	eventAddr string
	eventSock *zmq4.Socket
	commandAddr string
	commandSock *zmq4.Socket
}

func (g Guise) listenAndApply() {
	g.commandSock.Connect(g.commandAddr)
	subs := make([]SubCommand, 0)
	for {
		someBytes, _ := g.commandSock.RecvBytes(0)
		// os.Stderr.WriteString("guise: " + string(someBytes) + "\n")
		v := parseCommand(someBytes)
		switch cmd := v.(type) {
		case SubCommand: // odd: the go driver is collecting subs...
			subs = append(subs, cmd)
			cmd.Apply(g.View)
		case PostHTMLCommand:
			cmd.Apply(g.View)
			for _, sub := range subs {
				sub.Apply(g.View)
			}
		default: // odd: ... but the webview is collecting CSS
			cmd.Apply(g.View)
		}
	}
}

func NewGuise() Guise {
	eventAddr := "ipc:///tmp/guiseEvents"
	commandAddr := "ipc:///tmp/guiseCommands"
	fmt.Printf(`{"events":"%s", "commands":"%s"}`, eventAddr, commandAddr)
	os.Stdout.Close() // hmm, this seems like a bad global effect...
	eventSock, _ := zmq4.NewSocket(zmq4.PUSH)
	eventSock.Bind(eventAddr)
	eventSock.Send(`["hi"]`, 0)
	commandSock, _ := zmq4.NewSocket(zmq4.PULL)
	
	handleRPC := func(w webview.WebView, data string) {
		// os.Stderr.WriteString("guise-out: " + data + "\n")
		eventSock.Send(data, 0)
	}

	view := webview.New(webview.Settings{
		Width:     300,
		Height:    400,
		Title:     "Hi Guise",
		Resizable: true,
		ExternalInvokeCallback: handleRPC,
	})
	return Guise{view, eventAddr, eventSock, commandAddr, commandSock}
}

func main() {
	g := NewGuise()
	go g.listenAndApply()
	defer g.View.Exit()
	defer g.eventSock.Send(`["bye"]`, 0)
	
	// ???
	defer g.commandSock.Disconnect(g.commandAddr)
	defer g.eventSock.Disconnect(g.eventAddr)
	defer g.eventSock.Close()
	
	g.View.Run()
}