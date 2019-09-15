package main

import (
	"github.com/zserge/webview"
	"flag"
)


type Stui struct {
	View webview.WebView
	Conn
}

func (s Stui) listenAndApply() {
	subs := make([]SubscribeCommand, 0)
	for {
		v := s.Conn.Recv()
		switch cmd := v.(type) {
		case SubscribeCommand: // odd: the go driver is collecting subs...
			subs = append(subs, cmd)
			cmd.Apply(s.View)
		case PostElemCommand:
			cmd.Apply(s.View)
			for _, sub := range subs {
				sub.Apply(s.View)
			}
		case CloseGUICommand:
			cmd.Apply(s.View)
			break
		case CloseConnCommand:
			break
		default: // odd: ... but the webview is collecting CSS
			cmd.Apply(s.View)
		}
	}
}

func NewStui(c Conn) Stui {
	c.Start()
		
	cb := func(w webview.WebView, s string) {c.Send(s)}

	view := webview.New(webview.Settings{
		Width:     300,
		Height:    400,
		Title:     "Hi Stui",
		Resizable: true,
		ExternalInvokeCallback: cb,
	})
	c.Send(`["hi"]`)
	return Stui{view, c}
}

func main() {
	zmq := flag.String("zmq", "", "Socket name, if using zmq, such as ipc:///tmp/stui. Will use stdio if omitted or blank")
	flag.Parse()
	var conn Conn
	if *zmq == "" {
		conn = StdioConn()
	} else {
		conn = NewZMQConn(*zmq)
	}
	s := NewStui(conn)
	go s.listenAndApply()
	defer s.View.Exit()
	defer s.Send(`["bye"]`)
	defer s.Stop()
	
	s.View.Run()
}